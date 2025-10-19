package func008

import "context"

// Bad: context.Context is second parameter
func ProcessDataBad(id int, ctx context.Context) { // want "KTN-FUNC-008"
}

// Bad: context.Context is last parameter
func HandleRequestBad(id int, name string, ctx context.Context) { // want "KTN-FUNC-008"
}

// Bad: context in the middle
func ComplexFunctionBad(id int, ctx context.Context, name string) { // want "KTN-FUNC-008"
}

// Bad: Method with context not first
type ServiceBad struct{}

func (s *ServiceBad) Process(data string, ctx context.Context) { // want "KTN-FUNC-008"
}
