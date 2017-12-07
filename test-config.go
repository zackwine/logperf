package main

import (
  "fmt"
  "io/ioutil"
  "gopkg.in/yaml.v2"
)


type TestConfig struct {
  Url string
  Proxy string
  Count int
  Period int
  Message string
}

type TestConfigs struct {
    TestCfgs []TestConfig `tests`
}


func NewTestConfigs(filename string) (*TestConfigs, error) {
  t := &TestConfigs{}
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

