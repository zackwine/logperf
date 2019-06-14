package loggen

import (
  "log"
  "time"

  uuid "github.com/satori/go.uuid"
  "github.com/winez/logperf/outputs"
)

// FlowState the current state of the log flow
type FlowState string

const (
  // Init flow state is the state prior to sending logs
  Init FlowState = "Init"
  // Running flow state is when the log flow is sending logs
  Running FlowState = "Running"
  // Complete flow state is when the log flow has completed
  Complete FlowState = "Complete"
  // Failed flow state is when the log flow has failed
  Failed FlowState = "Failed"
)

// LogFlow : Create a log flow between a generator and output
type LogFlow struct {
  count      int
  loggen     *LogGenerator
  output     outputs.Output
  msgchan    chan string
  quittimer  chan bool
  log        *log.Logger
  TargetRate float64
  Sent       int
  StartTime  time.Time
  Elapsed    time.Duration
  UUID       string
  State      FlowState
}

// NewLogFlow : Initialize LogFlow
func NewLogFlow(output outputs.Output, component string, msgpadding int, daysoffset int, logger *log.Logger) *LogFlow {
  l := &LogFlow{}
  l.UUID = uuid.Must(uuid.NewV4()).String()
  l.loggen = NewLogGenerator(component, l.UUID)
  l.loggen.SetMessagePaddingSizeBytes(msgpadding)
  l.loggen.SetTimestampOffsetDays(daysoffset)
  l.output = output
  l.log = logger
  l.msgchan = make(chan string)
  l.quittimer = make(chan bool)
  l.State = Init
  return l
}

func (l *LogFlow) timeTrack(start time.Time, name string) {
  l.Elapsed = time.Since(start)
  messagesPerSec := l.GetMsgRate()
  log.Printf("%s took %s to write %d messages (%f per second, with target %f).", name, l.Elapsed, l.Sent, messagesPerSec, l.TargetRate)
}

// GetMsgRate - get the rate messages were sent
func (l *LogFlow) GetMsgRate() float64 {
  var elapsed time.Duration
  if l.State == Init || l.State == Failed {
    return 0
  }
  if l.State == Running {
    elapsed = time.Since(l.StartTime)
  } else if l.State == Complete {
    elapsed = l.Elapsed
  }

  elaspedSecs := float64(elapsed) / float64(time.Second)
  messagesPerSec := float64(l.Sent) / elaspedSecs
  return messagesPerSec
}

// The duration is only valid down to a value of about: 40*time.Microsecond
func (l *LogFlow) timerTask(period time.Duration, count int) error {
  l.TargetRate = float64(time.Second) / float64(period)
  ticker := time.NewTicker(period)
  err := l.output.StartOutput(l.msgchan)
  if err != nil {
    l.log.Printf("error: %v", err)
    return err
  }

  l.StartTime = time.Now()
  l.Elapsed = 0
  l.Sent = 0

  defer ticker.Stop()
  defer l.timeTrack(l.StartTime, "timerTask")

  for i := 0; i < count; i++ {
    select {
    case t := <-ticker.C:
      msg, err := l.loggen.GetMessage(t)
      if err != nil {
        l.log.Printf("error: %v", err)
      }
      l.msgchan <- msg

      l.Sent++
    case <-l.quittimer:
      l.log.Println("timerTask stopping")
      return nil
    }
  }

  l.log.Println("timerTask stopping based on count")

  return nil
}

func (l *LogFlow) stopTimerTask() {
  go func() {
    l.quittimer <- true
  }()
  err := l.output.StopOutput()
  if err != nil {
    l.log.Println("Failed to stop output plugin.")
  }
}

// Start : Start the log flow post to finished channel when task is complete
func (l *LogFlow) Start(period time.Duration, count int, finished chan *LogFlow) {
  go func() {
    err := l.timerTask(period, count)
    if err != nil {
      l.log.Println("Failed start timer task.")
    }
    finished <- l
  }()
}
