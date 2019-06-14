package logperf

import (
    "log"

    "time"

    "github.com/winez/logperf/loggen"
    "github.com/winez/logperf/outputs"
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
    } else {
        l.log.Fatalf("Invalid output specified %s", l.Output)
    }

    return loggen.NewLogFlow(output, l.Fields, l.Timefield, l.CounterField, l.Daysoffset, l.log)
}

// Start : Start a log perf test
// finished channel can be nil for blocking calls
func (l *LogPerf) Start(finished chan *LogPerf) error {

    // For each routine start a log flow
    for i := 0; i < l.Routines; i++ {
        lfr := l.initLogFlow()
        l.logFlowRoutines = append(l.logFlowRoutines, lfr)
        lfr.Start(time.Duration(l.Period)*time.Microsecond, l.Count, l.finishchan)
    }

    // Wait on all routines to complete
    for j := 0; j < l.Routines; j++ {
        select {
        case lf := <-l.finishchan:
            l.log.Printf("Finshed a flow with rate %f.", lf.GetMsgRate())
        }
    }

    // Notify the finished channel this completed if provided
    if finished != nil {
        finished <- l
    }

    return nil
}

// Stop : Start a log perf test
func (l *LogPerf) Stop() error {

    return nil
}
