package outputs

import (
	"errors"
	"fmt"
	"log"
)

// StdOutput : a channel based output for perftesting logs to a stdout input.
type StdOutput struct {
	log       *log.Logger
	inputchan chan string
	stopchan  chan bool
	running   bool
}

// NewStdOutput : Initialize StdOutput
func NewStdOutput(logger *log.Logger) *StdOutput {

	s := &StdOutput{
		log:      logger,
		stopchan: make(chan bool),
	}
	return s
}

// StartOutput : Implement the Output interface.
// Start reading from the channel and sending to the output.
func (s *StdOutput) StartOutput(input chan string) error {
	s.inputchan = input

	if s.running {
		log.Println("This output is already running")
		return errors.New("This output is already running")
	}
	s.running = true

	go func() {

		for {
			select {
			case message := <-s.inputchan:
				fmt.Println(message)
			case <-s.stopchan:
				s.log.Println("StdOutput thread stopping")
				s.running = false
				return
			}
		}

	}()

	return nil
}

// StopOutput : Implement the Output interface
// Stop reading from the channel and sending to the output.
func (s *StdOutput) StopOutput() error {
	if !s.running {
		log.Println("This output is NOT running")
		return errors.New("This output is NOT running")
	}
	go func() {
		s.stopchan <- true
	}()

	return nil
}
