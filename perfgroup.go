package logperf

import (
	"log"
)

// PerfGroup : Run a group of logperfs concurrently
type PerfGroup struct {
	logperfs   []*LogPerf
	log        *log.Logger
	finishchan chan *LogPerf
}

// NewPerfGroup : Initialize a group of logperf instances
func NewPerfGroup(cfgs []Config, logger *log.Logger) *PerfGroup {

	var logperfs []*LogPerf

	for _, conf := range cfgs {
		curPerf := NewLogPerf(conf, logger)
		logperfs = append(logperfs, curPerf)
	}

	return &PerfGroup{
		logperfs:   logperfs,
		log:        logger,
		finishchan: make(chan *LogPerf),
	}
}

// Start : Start a group of concurrent log perf tests
func (l *PerfGroup) Start() error {

	for _, logperf := range l.logperfs {
		err := logperf.Start(l.finishchan)
		if err != nil {
			l.log.Print(err)
			return err
		}
	}

	// Wait on all perfs to complete
	for _ = range l.logperfs {
		select {
		case lp := <-l.finishchan:
			l.log.Printf("Finshed a logperf test %s", lp.Name)
		}
	}

	return nil
}

// Stop : Start a log perf test
func (l *PerfGroup) Stop() error {
	for _, lp := range l.logperfs {
		lp.Stop()
	}

	return nil
}
