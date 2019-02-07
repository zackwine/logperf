package main

import (
  "log"
  "time"

  uuid "github.com/satori/go.uuid"
)

// LogFlow : Create a log flow between a generator and output
type LogFlow struct {
  count      int
  loggen     *LogGenerator
  output     Output
  msgchan    chan string
  quittimer  chan bool
  log        *log.Logger
  TargetRate float64
  Sent       int
  StartTime  time.Time
  Elapsed    time.Duration
  UUID       string
}

// NewLogFlow : Initialize LogFlow
func NewLogFlow(output Output, logger *log.Logger) *LogFlow {
  l := &LogFlow{}
  l.UUID = uuid.Must(uuid.NewV4()).String()
  l.loggen = NewLogGenerator("LogFlow", l.UUID)
  l.loggen.SetMessagePaddingSizeBytes(300)
  l.output = output
  l.log = logger
  l.msgchan = make(chan string)
  l.quittimer = make(chan bool)
  return l
}

func (l *LogFlow) timeTrack(start time.Time, name string) {
  l.Elapsed = time.Since(start)
  messagesPerSec := l.getMsgRate()
  log.Printf("%s took %s to write %d messages (%f per second, with target %f).", name, l.Elapsed, l.Sent, messagesPerSec, l.TargetRate)
}

func (l *LogFlow) getMsgRate() float64 {
  var elapsed time.Duration
  if l.Elapsed == 0 {
    elapsed = l.Elapsed
  } else {
    elapsed = time.Since(l.StartTime)
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
    logger.Printf("error: %v", err)
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
        logger.Printf("error: %v", err)
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
  l.output.StopOutput()
}
