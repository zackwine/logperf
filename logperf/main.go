package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/zackwine/logperf"
	"github.com/zackwine/logperf/api"
)

var (
	perffile   = flag.String("perffile", "", "The path to the logperf file (yaml) defining perf tests to run.")
	http       = flag.Bool("http", false, "Enable restful API to start perf tests.")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	logfile    = flag.Bool("logfile", false, "Log to file /var/log/logperf.log")
)

func printUsageAndBail(message string) {
	fmt.Fprintln(os.Stderr, "ERROR:", message)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Usage:")
	flag.PrintDefaults()
	os.Exit(64)
}

func runperffile(logger *log.Logger) {

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

func main() {

	flag.Parse()

	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	if *logfile {
		f, err := os.OpenFile("/var/log/logperf.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			logger.Fatalln(err)
		}
		logger = log.New(f, "", log.LstdFlags|log.Lshortfile)
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
		httpServer := &api.HTTPServer{Logger: logger, Addr: ":8080"}
		// Server runs async
		httpServer.Start()
		// Wait on signal from user
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigs
		httpServer.Stop()
		fmt.Println("Exiting with signal:")
		fmt.Println(sig)

	} else if *perffile != "" {
		runperffile(logger)
	} else {
		flag.PrintDefaults()
	}

}
