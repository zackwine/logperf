package main

import (
  "log"
  "time"
)

const (
    maxConnectsPerHost = 25
    requestTimeout = 5
)

type LogPerf struct {
  count int
  loggen *LogGenerator
  httpoutput *HttpOutputPool
}


func NewLogPerf(count int) (*LogPerf) {
  l := &LogPerf{}
  l.count = count
  l.loggen = NewLogGenerator("logperf")
  l.loggen.SetMessagePaddingSizeBytes(300)
  l.httpoutput = NewHttpOutputPool( maxConnectsPerHost, requestTimeout)
  return l
}

func timeTrack(start time.Time, name string) {
    elapsed := time.Since(start)
    log.Printf("%s took %s", name, elapsed)
}

func (l *LogPerf) sendLog() (error) {
  url := "http://64.102.245.5:14001/logs.dp.vmr"
  //url := "http://op-ingest-fluentd.elb.load-2.ciscoanalytics.com:24224/logs.dp.vmr"
  msg, err := l.loggen.GetMessage()
  if err != nil {
    logger.Printf("error: %v", err)
    return err
  }
  defer timeTrack(time.Now(), "SendMessage")
  l.httpoutput.SendMessage(url, "POST", "json=" + msg)
  return nil
}

func (l *LogPerf) SendLogs() {
  for i := 0; i < l.count; i++ {
    l.sendLog()
    //time.Sleep(4 * time.Second)
  }
}