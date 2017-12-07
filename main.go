package main

import (
  "flag"
  "fmt"
  "log"
  "os"
  "runtime/pprof"
)

var (
  testFile    = flag.String("testfile", "test.yaml", "The path to the test file (yaml) defining tests to run.")
  cpuprofile  = flag.String("cpuprofile", "", "write cpu profile to file")

  logger = log.New(os.Stderr, "", log.LstdFlags)
)


func printUsageAndBail(message string) {
  fmt.Fprintln(os.Stderr, "ERROR:", message)
  fmt.Fprintln(os.Stderr)
  fmt.Fprintln(os.Stderr, "Usage:")
  flag.PrintDefaults()
  os.Exit(64)
}

func main() {
  flag.Parse()

  if *testFile == "" {
    printUsageAndBail("No test definition provided.")
  }

  if *cpuprofile != "" {
    f, err := os.Create(*cpuprofile)
    if err != nil {
      logger.Fatal(err)
    }
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
  }

  testConfigs, err := NewTestConfigs(*testFile)

  logger.Printf("error: %v", testConfigs)
  if err != nil {
    logger.Printf("error: %v", err)
  }else{
    logger.Printf("testConfigs: %v", testConfigs)
  }

  logperf := NewLogPerf(300)
  logperf.SendLogs()

}

