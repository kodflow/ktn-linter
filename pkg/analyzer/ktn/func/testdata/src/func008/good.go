package func008

import "context"

// Good: context.Context as first parameter
func ProcessWithContext(ctx context.Context, data string) {
}

// Good: context as first parameter with multiple params
func HandleRequestGood(ctx context.Context, id int, name string) {
}

// Good: No context parameter at all
func SimpleFunctionGood(id int, name string) {
}

// Good: Method with context as first parameter (after receiver)
type ServiceGood struct{}

func (s *ServiceGood) ProcessData(ctx context.Context, data string) {
}

// Good: No parameters
func NoParams() {
}

// Good: Context first with only one parameter
func OnlyContextParam(ctx context.Context) {
}

// Good: Using context.Background() - not a type but tests context package usage
func UseContextBackground() {
	_ = context.Background()
}

// Good: Function that takes context.CancelFunc (different from context.Context)
func TakesCancelFunc(id int, cancel context.CancelFunc) {
}

// Good: Test functions should be ignored even if context is not first
func TestSomething(t interface{}, ctx context.Context) {
}

func TestAnotherThing(data string, ctx context.Context) {
}

func BenchmarkProcess(b interface{}, ctx context.Context) {
}
