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
