package func002

// Bad: Function with 6 parameters
func SixParams(a, b, c, d, e, f int) { // want "KTN-FUNC-002"
}

// Bad: Function with 7 parameters of different types
func SevenParams(a int, b string, c bool, d float64, e []int, f map[string]int, g chan int) { // want "KTN-FUNC-002"
}

// Bad: Function with many parameters
func TooManyParams(a, b, c, d, e, f, g, h, i, j int) { // want "KTN-FUNC-002"
}

// Bad: Function literal with too many params
var BadLiteral = func(a, b, c, d, e, f int) { // want "KTN-FUNC-002"
}

// Bad: Function literal with unnamed parameters (more than 5)
var BadLiteralUnnamed = func(int, string, bool, float64, []int, map[string]int) { // want "KTN-FUNC-002"
}

// Bad: Function with mixed named and unnamed is not possible in Go,
// but we test unnamed params in func literals
var BadLiteralSixUnnamed = func(int, int, int, int, int, int) { // want "KTN-FUNC-002"
}

// Bad: Function with 8 parameters (edge case for counting)
func EightParams(a int, b, c string, d, e, f bool, g, h float64) { // want "KTN-FUNC-002"
}

// Bad: Function with exactly 6 parameters (just over limit)
func ExactlySixParams(a, b, c, d, e, f int) { // want "KTN-FUNC-002"
}

// Bad: Function with variadic parameter but still over limit (6 params)
func WithVariadicBad(a, b, c, d, e int, f ...string) { // want "KTN-FUNC-002"
}
