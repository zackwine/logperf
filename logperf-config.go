package main

import (
  "fmt"
  "io/ioutil"

  "gopkg.in/yaml.v2"
)

// LogPerfConfig : Configruation to run log performance tests
type LogPerfConfig struct {
  Output     string
  Addr       string
  Count      int
  Period     int
  Daysoffset int
  Padding    int
  Message    string
  Component  string
  Routines   int
}

// LogPerfConfigs : Configruation to run log performance tests
type LogPerfConfigs struct {
  LogPerfConfigs []LogPerfConfig `logperfs`
}

// NewLogPerfConfigs : Unmarshal LogPerfConfig from YAML file
func NewLogPerfConfigs(filename string) (*LogPerfConfigs, error) {
  t := &LogPerfConfigs{}
  source, err := ioutil.ReadFile(filename)
  if err != nil {
    panic(err)
  }
  err = yaml.Unmarshal(source, t)
  if err != nil {
    fmt.Printf("error: %v", err)
    panic(err)
  }
  return t, err
}
