package logperf

import (
	"log"

	"time"

	"github.com/zackwine/logperf/loggen"
	"github.com/zackwine/logperf/outputs"
)

// LogPerf : a channel based output for perftesting logs to a TCP input.
type LogPerf struct {
	Config
	log             *log.Logger
	finishchan      chan *loggen.LogFlow
	logFlowRoutines []*loggen.LogFlow
}

// NewLogPerf : Initialize LogPerf
func NewLogPerf(conf Config, logger *log.Logger) *LogPerf {

	return &LogPerf{
		Config:     conf,
		log:        logger,
		finishchan: make(chan *loggen.LogFlow),
	}
}

func (l *LogPerf) initLogFlow() *loggen.LogFlow {
	var output outputs.Output
	l.log.Printf("Configuring output (%s)...", l.Output)
	if l.Output == "tcp" {
		output = outputs.NewTCPOutput(l.Addr, l.log)
	} else if l.Output == "stdout" {
		output = outputs.NewStdOutput(l.log)
	} else if l.Output == "elasticsearch" {
		esConfig := outputs.ElasticSearchOutputConfig{
			Addresses: []string{l.Addr},
			IndexBase: l.IndexBase,
		}
		output = outputs.NewElasticSearchOutput(esConfig, l.log)
	} else {
		l.log.Fatalf("Invalid output specified %s", l.Output)
	}

	return loggen.NewLogFlow(output, l.OutputFormat, l.Fields, l.Timefield, l.CounterField, l.Daysoffset, l.log)
}

// GetTargetCount - Get the target number of logs to be generated by all routines
func (l *LogPerf) GetTargetCount() int64 {
	return l.Count * int64(l.Routines)
}

// GetCurrentCount - Get the number of logs generated thus far
func (l *LogPerf) GetCurrentCount() int64 {
	var curCnt int64
	for _, lfr := range l.logFlowRoutines {
		curCnt += lfr.Sent
	}
	return curCnt
}

// Start : Start a log perf test
// finchan will be notified when all routines have completed
func (l *LogPerf) Start(finchan chan *LogPerf) error {

	// For each routine start a log flow
	for i := 0; i < l.Routines; i++ {
		lfr := l.initLogFlow()
		l.logFlowRoutines = append(l.logFlowRoutines, lfr)
		lfr.Start(time.Duration(l.Period)*time.Microsecond, l.Count, l.finishchan)
	}

	go func() {
		// Wait on all routines to complete
		for j := 0; j < l.Routines; j++ {
			select {
			case lf := <-l.finishchan:
				l.log.Printf("Finshed a flow with rate %f.", lf.GetMsgRate())
			}
		}

		// Notify the finished channel this completed if provided
		if finchan != nil {
			finchan <- l
		}
	}()

	return nil
}

// Stop : Start a log perf test
func (l *LogPerf) Stop() {
	// Stop all log flows
	for _, lfr := range l.logFlowRoutines {
		lfr.Stop()
	}
}
