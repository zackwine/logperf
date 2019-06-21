package loggen

import (
  "log"
  "time"

  "github.com/zackwine/logperf/outputs"
)

// FlowState the current state of the log flow
type FlowState string

const (
  // Init flow state is the state prior to sending logs
  Init FlowState = "Init"
  // Running flow state is when the log flow is sending logs
  Running FlowState = "Running"
  // Completed flow state is when the log flow has completed
  Completed FlowState = "Completed"
  // Stopped flow state is when the log flow has been stopped prematurely
  Stopped FlowState = "Stopped"
)

// LogFlow : Create a log flow between a generator and output
type LogFlow struct {
  Count      int64
  loggen     *LogGenerator
  output     outputs.Output
  msgchan    chan string
  quittimer  chan bool
  log        *log.Logger
  TargetRate float64
  Sent       int
  StartTime  time.Time
  Elapsed    time.Duration
  State      FlowState
}

// NewLogFlow : Initialize LogFlow
func NewLogFlow(output outputs.Output, fields map[string]interface{}, timeField string, counterField string, daysoffset int, logger *log.Logger) *LogFlow {
  l := &LogFlow{}
  l.loggen = NewLogGenerator(fields, timeField, counterField, daysoffset, logger)
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
  if l.State == Init {
    return 0
  }
  if l.State == Running {
    elapsed = time.Since(l.StartTime)
  } else if l.State == Completed {
    elapsed = l.Elapsed
  }

  elaspedSecs := float64(elapsed) / float64(time.Second)
  messagesPerSec := float64(l.Sent) / elaspedSecs
  return messagesPerSec
}

// The duration is only valid down to a value of about: 40*time.Microsecond
func (l *LogFlow) timerTask(period time.Duration, count int64) error {
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
  l.State = Running

  defer ticker.Stop()
  defer l.timeTrack(l.StartTime, "timerTask")

  for i := int64(0); i < count; i++ {
    select {
    case t := <-ticker.C:
      msg, err := l.loggen.GetMessage(t)
      if err != nil {
        l.log.Printf("error: %v", err)
      }
      l.msgchan <- msg

      l.Sent++
    case <-l.quittimer:
      l.State = Stopped
      l.log.Println("timerTask stopping")
      return nil
    }
  }

  l.log.Println("timerTask stopping based on count")
  l.State = Completed

  return nil
}

// Stop : Stop the this log flow
func (l *LogFlow) Stop() {
  if l.State != Running {
    l.log.Println("LogFlow not running.")
    return
  }
  go func() {
    l.quittimer <- true
  }()
  err := l.output.StopOutput()
  if err != nil {
    l.log.Println("Failed to stop output plugin.")
  }
}

// Start : Start the log flow post to finished channel when task is complete
func (l *LogFlow) Start(period time.Duration, count int64, finished chan *LogFlow) {
  go func() {
    l.Count = count
    err := l.timerTask(period, count)
    if err != nil {
      l.log.Println("Failed start timer task.")
    }
    finished <- l
  }()
}
