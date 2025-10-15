package KTN_INTERFACE_005

type processor struct{}

func (p *processor) Process(data string) string {
    return data
}

func newProcessor() *processor {
    return &processor{}
}
