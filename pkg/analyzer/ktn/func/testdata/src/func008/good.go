package func008

import "context"

// ProcessWithContext demonstrates context.Context as first parameter.
//
// Params:
//   - ctx: context for cancellation
//   - data: input data string
func ProcessWithContext(ctx context.Context, data string) {
}

// HandleRequestGood demonstrates context as first parameter with multiple params.
//
// Params:
//   - ctx: context for cancellation
//   - id: identifier
//   - name: name string
func HandleRequestGood(ctx context.Context, id int, name string) {
}

// SimpleFunctionGood demonstrates no context parameter at all.
//
// Params:
//   - id: identifier
//   - name: name string
func SimpleFunctionGood(id int, name string) {
}

// ServiceGood demonstrates a service type.
type ServiceGood struct{}

// ProcessData demonstrates method with context as first parameter after receiver.
//
// Params:
//   - ctx: context for cancellation
//   - data: input data string
func (s *ServiceGood) ProcessData(ctx context.Context, data string) {
}

// NoParams demonstrates a function with no parameters.
func NoParams() {
}

// OnlyContextParam demonstrates context first with only one parameter.
//
// Params:
//   - ctx: context for cancellation
func OnlyContextParam(ctx context.Context) {
}

// UseContextBackground demonstrates using context.Background() - not a type but tests context package usage.
func UseContextBackground() {
	_ = context.Background()
}

// TakesCancelFunc demonstrates a function that takes context.CancelFunc (different from context.Context).
//
// Params:
//   - id: identifier
//   - cancel: cancellation function
func TakesCancelFunc(id int, cancel context.CancelFunc) {
}

// TestSomething demonstrates test functions should be ignored even if context is not first.
//
// Params:
//   - t: testing interface
//   - ctx: context for cancellation
func TestSomething(t interface{}, ctx context.Context) {
}

// TestAnotherThing demonstrates another test function.
//
// Params:
//   - data: test data string
//   - ctx: context for cancellation
func TestAnotherThing(data string, ctx context.Context) {
}

// BenchmarkProcess demonstrates a benchmark function.
//
// Params:
//   - b: benchmark interface
//   - ctx: context for cancellation
func BenchmarkProcess(b interface{}, ctx context.Context) {
}
