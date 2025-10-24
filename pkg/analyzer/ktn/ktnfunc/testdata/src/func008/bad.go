package func008

import "context"

// ProcessDataBad demonstrates context.Context as second parameter (intentional violation).
//
// Params:
//   - id: identifier
//   - ctx: context for cancellation (BAD: should be first)
func ProcessDataBad(id int, ctx context.Context) {
}

// HandleRequestBad demonstrates context.Context as last parameter (intentional violation).
//
// Params:
//   - id: identifier
//   - name: name string
//   - ctx: context for cancellation (BAD: should be first)
func HandleRequestBad(id int, name string, ctx context.Context) {
}

// ComplexFunctionBad demonstrates context.Context in the middle (intentional violation).
//
// Params:
//   - id: identifier
//   - ctx: context for cancellation (BAD: should be first)
//   - name: name string
func ComplexFunctionBad(id int, ctx context.Context, name string) {
}

// ServiceBad demonstrates a service type with bad context positioning.
type ServiceBad struct{}

// Process demonstrates method with context not first (intentional violation).
//
// Params:
//   - data: input data string
//   - ctx: context for cancellation (BAD: should be first)
func (s *ServiceBad) Process(data string, ctx context.Context) {
}
