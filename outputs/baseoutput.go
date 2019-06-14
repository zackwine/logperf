package outputs

// Output : Interface for loggen to sink logs towards an output.
type Output interface {
	StartOutput(input chan string) error
	StopOutput() error
}
