package const002

// File with mix of const, var, and function declarations
const (
	WithFunc1 string = "func1"
	WithFunc2 string = "func2"
)

// This function is a non-GenDecl that should be skipped
func SomeFunction() string {
	return "test"
}

var WithFuncVar string = "var"
