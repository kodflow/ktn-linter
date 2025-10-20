package func008

import "context"

// Bad: context.Context is second parameter
func ProcessDataBad(id int, ctx context.Context) {
}

// Bad: context.Context is last parameter
func HandleRequestBad(id int, name string, ctx context.Context) {
}

// Bad: context in the middle
func ComplexFunctionBad(id int, ctx context.Context, name string) {
}

// Bad: Method with context not first
type ServiceBad struct{}

func (s *ServiceBad) Process(data string, ctx context.Context) {
}
