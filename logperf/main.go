package main

import (
  "flag"
  "fmt"
  "log"
  "os"
  "runtime/pprof"

  "github.com/zackwine/logperf"
  "github.com/zackwine/logperf/api"
)

var (
  perffile   = flag.String("perffile", "logperf.yaml", "The path to the logperf file (yaml) defining perf tests to run.")
  http       = flag.Bool("http", false, "Enable restful API to start perf tests.")
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
  logger.SetFlags(log.LstdFlags | log.Lshortfile)
  flag.Parse()

  if *perffile == "" {
    printUsageAndBail("No test definition provided.")
  }

  if *cpuprofile != "" {
    f, err := os.Create(*cpuprofile)
    if err != nil {
      logger.Fatal(err)
    }
    err = pprof.StartCPUProfile(f)
    if err != nil {
      logger.Fatal(err)
    }
    defer pprof.StopCPUProfile()
  }

  if *http {
    api.RunServer(logger)
  }

  cfgs, err := logperf.NewConfigs(*perffile)
  if err != nil {
    logger.Printf("error: %v", err)
  } else {
    logger.Printf("perfConfigs: %v", cfgs)
  }

  logperf := logperf.NewPerfGroup(cfgs.Configs, logger)

  err = logperf.Start()
  if err != nil {
    logger.Fatal(err)
  }

}
