package const003

import "fmt"

// Type declarations should be skipped by const003
type MyStruct struct {
	Field string
}

type MyInterface interface {
	Method()
}

// Valid const
const VALID_CONST = "test"

var testVar = fmt.Sprintf("var")
