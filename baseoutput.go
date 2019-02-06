package main

type Output interface {
	StartOutput(input chan string) error
	StopOutput() error
}
