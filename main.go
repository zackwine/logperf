package main

import (
  "flag"
  "fmt"
  "log"
  "os"
  "runtime/pprof"
)

var (
  perffile   = flag.String("perffile", "logperf.yaml", "The path to the logperf file (yaml) defining perf tests to run.")
  cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

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

  if *perffile == "" {
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

  perfConfigs, err := NewLogPerfConfigs(*perffile)
  if err != nil {
    logger.Printf("error: %v", err)
  } else {
    logger.Printf("perfConfigs: %v", perfConfigs)
  }

  logperf := NewLogPerf(perfConfigs.LogPerfConfigs, logger)

  logperf.Start()

  //tcp := NewTCPOutput("127.0.0.1:5000", logger)
  //tcp := NewTCPOutput("lmm-logstash-collector:4500", logger)
  //logflow := NewLogFlow(tcp, "logperf", 300, 0, logger)
  //go logflow.timerTask(30*time.Microsecond, 100000)

  //tcp2 := NewTCPOutput("lmm-logstash-collector:4500", logger)
  //logflow2 := NewLogFlow(tcp2, "logperf2", 200, 0, logger)
  //logflow2.timerTask(30*time.Microsecond, 100000)
}
