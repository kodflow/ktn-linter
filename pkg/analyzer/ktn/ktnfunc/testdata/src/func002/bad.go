// Bad examples for the func002 test case.
package func002

import "context"

// ProcessDataBad demonstrates context.Context as second parameter (intentional violation).
//
// Params:
//   - id: identifier
//   - ctx: context for cancellation (BAD: should be first)
func ProcessDataBad(id int, ctx context.Context) {
	// Utilisation des paramètres
	_, _ = id, ctx
}

// HandleRequestBad demonstrates context.Context as last parameter (intentional violation).
//
// Params:
//   - id: identifier
//   - name: name string
//   - ctx: context for cancellation (BAD: should be first)
func HandleRequestBad(id int, name string, ctx context.Context) {
	// Utilisation des paramètres
	_, _, _ = id, name, ctx
}

// ComplexFunctionBad demonstrates context.Context in the middle (intentional violation).
//
// Params:
//   - id: identifier
//   - ctx: context for cancellation (BAD: should be first)
//   - name: name string
func ComplexFunctionBad(id int, ctx context.Context, name string) {
	// Utilisation des paramètres
	_, _, _ = id, ctx, name
}

// ServiceBad demonstrates a service type with bad context positioning.
// Used to test context position in method receivers.
type ServiceBad struct{}

// ServiceBadInterface defines the public methods of ServiceBad.
type ServiceBadInterface interface {
	Process(data string, ctx context.Context)
}

// NewServiceBad creates a new instance of ServiceBad.
//
// Returns:
//   - *ServiceBad: new instance
func NewServiceBad() *ServiceBad {
	// Return new instance
	return &ServiceBad{}
}

// Process demonstrates method with context not first (intentional violation).
//
// Params:
//   - data: input data string
//   - ctx: context for cancellation (BAD: should be first)
func (s *ServiceBad) Process(data string, ctx context.Context) {
	// Utilisation des paramètres
	_, _ = data, ctx
}
