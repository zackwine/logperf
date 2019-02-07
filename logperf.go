package main

import (
    "log"

    "time"
)

type logFlowWrapper struct {
    logPerfConfig   LogPerfConfig
    logFlowRoutines []*LogFlow
}

// LogPerf : a channel based output for perftesting logs to a TCP input.
type LogPerf struct {
    flows      []*logFlowWrapper
    log        *log.Logger
    finishchan chan *LogFlow
}

// NewLogPerf : Initialize LogPerf
func NewLogPerf(logPerfConfigs []LogPerfConfig, logger *log.Logger) *LogPerf {

    var flows []*logFlowWrapper

    for _, conf := range logPerfConfigs {
        curFlow := &logFlowWrapper{
            logPerfConfig: conf,
        }
        flows = append(flows, curFlow)
    }

    l := &LogPerf{
        log:        logger,
        flows:      flows,
        finishchan: make(chan *LogFlow),
    }
    return l
}

func (l *LogPerf) initLogFlow(logPerfConfig LogPerfConfig) *LogFlow {
    var output Output
    if logPerfConfig.Output == "tcp" {
        output = NewTCPOutput(logPerfConfig.Addr, logger)
    } else {

    }
    return NewLogFlow(output, logPerfConfig.Component, logPerfConfig.Padding, logPerfConfig.Daysoffset, l.log)
}

// Start : Start a log perf test
func (l *LogPerf) Start() error {
    var flowCnt int
    for _, flow := range l.flows {
        // For each routine start a log flow
        for i := 0; i < flow.logPerfConfig.Routines; i++ {
            lfr := l.initLogFlow(flow.logPerfConfig)
            flow.logFlowRoutines = append(flow.logFlowRoutines, lfr)
            lfr.Start(time.Duration(flow.logPerfConfig.Period)*time.Microsecond, flow.logPerfConfig.Count, l.finishchan)
            flowCnt++
        }
    }

    for j := 0; j < flowCnt; j++ {
        select {
        case lf := <-l.finishchan:
            l.log.Printf("Finshed a flow with rate %f.", lf.getMsgRate())
        }
    }

    return nil
}

// Stop : Start a log perf test
func (l *LogPerf) Stop() error {

    return nil
}
