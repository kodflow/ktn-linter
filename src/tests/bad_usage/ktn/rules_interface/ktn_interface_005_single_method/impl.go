package KTN_INTERFACE_005

type processor struct{}

func (p *processor) Process(data string) string {
	// Early return from function.
	return data
}

func newProcessor() *processor {
	// Early return from function.
	return &processor{}
}
