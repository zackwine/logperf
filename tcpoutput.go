package main

import (
    "errors"
    "log"
    "net"
)

// TCPOutput : a channel based output for perftesting logs to a TCP input.
type TCPOutput struct {
    address   string
    conn      net.Conn
    log       *log.Logger
    inputchan chan string
    stopchan  chan bool
    running   bool
}

// NewTCPOutput : Initialize TCPOutput
func NewTCPOutput(address string, logger *log.Logger) *TCPOutput {

    h := &TCPOutput{
        address:  address,
        log:      logger,
        stopchan: make(chan bool),
    }
    return h
}

// StartOutput : Implement the Output interface.
// Start reading from the channel and sending to the output.
func (t *TCPOutput) StartOutput(input chan string) error {
    t.inputchan = input

    err := t.connect()
    if err != nil {
        return err
    }

    if t.running {
        log.Println("This output is already running")
        return errors.New("This output is already running")
    }
    t.running = true

    go func() {

        for {
            select {
            case message := <-t.inputchan:
                t.write(message + "\n")
                //t.log.Println("Writing message", message)
            case <-t.stopchan:
                t.log.Println("StartOutput stopping")
                t.running = false
                return
            }
        }

    }()

    return nil
}

// StopOutput : Implement the Output interface
// Stop reading from the channel and sending to the output.
func (t *TCPOutput) StopOutput() error {
    if !t.running {
        log.Println("This output is NOT running")
        return errors.New("This output is NOT running")
    }
    go func() {
        t.stopchan <- true
    }()
    t.close()

    return nil
}

func (t *TCPOutput) connect() error {
    var err error
    t.conn, err = net.Dial("tcp", t.address)
    if err != nil {
        t.log.Println("Failed to connect to", t.address)
        return err
    }
    return nil
}

func (t *TCPOutput) write(message string) error {
    // First var is an int bytesWritten?
    _, err := t.conn.Write([]byte(message + "\n"))
    if err != nil {
        t.log.Println("Failed to write to", t.address)
        return err
    }
    return nil
}

func (t *TCPOutput) close() error {
    err := t.conn.Close()
    if err != nil {
        t.log.Println("Failed to close", t.address)
        return err
    }
    return nil
}
