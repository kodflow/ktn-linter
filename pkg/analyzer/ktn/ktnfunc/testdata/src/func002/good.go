// Package func002 contient des exemples de fonctions utilisant context.Context.
package func002

import "context"

// ProcessWithContext demonstrates context.Context as first parameter.
//
// Params:
//   - _ctx: context for cancellation (non utilisé dans cet exemple)
//   - _data: input data string (non utilisé dans cet exemple)
func ProcessWithContext(_ctx context.Context, _data string) {
}

// HandleRequestGood demonstrates context as first parameter with multiple params.
//
// Params:
//   - _ctx: context for cancellation (non utilisé dans cet exemple)
//   - _id: identifier (non utilisé dans cet exemple)
//   - _name: name string (non utilisé dans cet exemple)
func HandleRequestGood(_ctx context.Context, _id int, _name string) {
}

// SimpleFunctionGood demonstrates no context parameter at all.
//
// Params:
//   - _id: identifier (non utilisé dans cet exemple)
//   - _name: name string (non utilisé dans cet exemple)
func SimpleFunctionGood(_id int, _name string) {
}

// ServiceGood demonstrates a service type.
// Provides data processing capabilities with context support.
type ServiceGood struct{}

// ServiceGoodInterface définit les méthodes publiques de ServiceGood.
type ServiceGoodInterface interface {
	ProcessData(ctx context.Context, data string)
}

// NewServiceGood crée une nouvelle instance de ServiceGood.
//
// Returns:
//   - *ServiceGood: nouvelle instance du service
func NewServiceGood() *ServiceGood {
	// Retour de la nouvelle instance
	return &ServiceGood{}
}

// ProcessData demonstrates method with context as first parameter after receiver.
//
// Params:
//   - _ctx: context for cancellation (non utilisé dans cet exemple)
//   - _data: input data string (non utilisé dans cet exemple)
func (s *ServiceGood) ProcessData(_ctx context.Context, _data string) {
}

// NoParams demonstrates a function with no parameters.
func NoParams() {
}

// OnlyContextParam demonstrates context first with only one parameter.
//
// Params:
//   - _ctx: context for cancellation (non utilisé dans cet exemple)
func OnlyContextParam(_ctx context.Context) {
}

// UseContextBackground demonstrates using context.Background() - not a type but tests context package usage.
func UseContextBackground() {
	_ = context.Background()
}

// TakesCancelFunc demonstrates a function that takes context.CancelFunc (different from context.Context).
//
// Params:
//   - _id: identifier (non utilisé dans cet exemple)
//   - _cancel: cancellation function (non utilisée dans cet exemple)
func TakesCancelFunc(_id int, _cancel context.CancelFunc) {
}

// TestSomething demonstrates test functions should be ignored even if context is not first.
//
// Params:
//   - _t: testing interface (non utilisé dans cet exemple)
//   - _ctx: context for cancellation (non utilisé dans cet exemple)
func TestSomething(_t any, _ctx context.Context) {
}

// TestAnotherThing demonstrates another test function.
//
// Params:
//   - _data: test data string (non utilisé dans cet exemple)
//   - _ctx: context for cancellation (non utilisé dans cet exemple)
func TestAnotherThing(_data string, _ctx context.Context) {
}

// BenchmarkProcess demonstrates a benchmark function.
//
// Params:
//   - _b: benchmark interface (non utilisé dans cet exemple)
//   - _ctx: context for cancellation (non utilisé dans cet exemple)
func BenchmarkProcess(_b any, _ctx context.Context) {
}
