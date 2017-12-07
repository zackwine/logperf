package main

import (
  "encoding/json"
  "strconv"
  "time"
)

type LogGenerator struct {
  baseMap map[string]string
  baseMessage string
  messagePaddingSizeBytes int
  timestampFieldName string
  timestampOffsetDays int
  seqNum int
  randstrgen *RandStringGen
}

func NewLogGenerator(componentName string) (*LogGenerator) {
  l := &LogGenerator{}
  // Set default timestamp field
  l.randstrgen = NewRandStringGen()
  l.timestampFieldName = "@timestamp"
  l.baseMessage = "Fake message: "
  l.baseMap = make(map[string]string)
  // Create some default fake JSON fields for the log
  l.baseMap["logger"] = "loggen-util-test"
  l.baseMap["nodeid"] = "0"
  l.baseMap["nodetype"] = "fake"
  l.baseMap["logsource"] = "console"
  l.baseMap["component"] = componentName
  l.baseMap["subsystem"] = "loggen"
  l.baseMap["ssinst"] = "007"
  l.baseMap["dc"] = "node.fake.vci"
  l.baseMap["host"] = componentName
  l.baseMap["_index"] = "logstash-2017-12-04"
  return l
}

func(l *LogGenerator) SetBaseMessage( baseMessage string ) {
  l.baseMessage = baseMessage
}

func(l *LogGenerator) SetMessagePaddingSizeBytes( messagePaddingSizeBytes int ) {
  l.messagePaddingSizeBytes = messagePaddingSizeBytes
}

func(l *LogGenerator) SetTimestampFieldName( timestampFieldName string ) {
  l.timestampFieldName = timestampFieldName
}

func(l *LogGenerator) SetTimestampOffsetDays( timestampOffsetDays int ) {
  l.timestampOffsetDays = timestampOffsetDays
}

func(l *LogGenerator) SetField( field string, value string ) {
  l.baseMap[field] = value
}

func(l *LogGenerator) ResetSeqNum() {
  l.seqNum = 0
}

func(l *LogGenerator) GetMessage( ) (string, error) {

  if l.messagePaddingSizeBytes > 0 {
    l.baseMap["message"] = l.baseMessage + " " + l.randstrgen.RandString(l.messagePaddingSizeBytes)
  }else{
    l.baseMap["message"] = l.baseMessage 
  }
  l.baseMap["seqNum"] = strconv.Itoa(l.seqNum)
  l.seqNum++
  offsetstr := time.Now().AddDate(0, 0, -l.timestampOffsetDays).Format(time.RFC3339)
  l.baseMap[l.timestampFieldName] = offsetstr

  messagebytes, err := json.Marshal(l.baseMap)
  if err != nil {
    logger.Println(err)
    return "", err
  }

  return string(messagebytes[:]), err
}

func(l *LogGenerator) GenMessages(count int, period int) {

}
