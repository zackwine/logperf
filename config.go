package logperf

import (
  "fmt"
  "io/ioutil"
  "path/filepath"

  "gopkg.in/yaml.v2"
)

// Config : Configruation to run log performance tests
type Config struct {
  Name       string `json:"name"`
  Output     string `json:"output"`
  Addr       string `json:"addr"`
  Count      int    `json:"count"`
  Period     int    `json:"period"`
  Daysoffset int    `json:"daysoffset"`
  Padding    int    `json:"padding"`
  Message    string `json:"message"`
  Component  string `json:"component"`
  Routines   int    `json:"routines"`
}

// Configs : A list of logperf configs that can be Marshalled from yaml
type Configs struct {
  Configs []Config `yaml:"logperfs"`
}

// NewConfigs : Unmarshal LogPerfConfig from YAML file
func NewConfigs(filename string) (*Configs, error) {
  t := &Configs{}
  source, err := ioutil.ReadFile(filepath.Clean(filename))
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
