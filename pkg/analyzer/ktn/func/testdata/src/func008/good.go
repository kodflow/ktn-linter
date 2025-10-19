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
