package func002

// Good: Function with no parameters
func NoParams() {
}

// Good: Function with 1 parameter
func OneParam(a int) {
}

// Good: Function with 5 parameters (exactly at limit)
func FiveParams(a, b, c, d, e int) {
}

// Good: Function with 5 parameters of different types
func FiveParamsMixed(a int, b string, c bool, d float64, e []int) {
}

// Good: Method with receiver + 4 parameters = 5 total params (receiver doesn't count)
type MyType struct{}

func (m MyType) MethodWithFourParams(a, b, c, d int) {
}

// Good: Test functions are exempt
func TestWithManyParams(t, a, b, c, d, e, f int) {
}

// Good: Benchmark functions are exempt
func BenchmarkWithManyParams(b, a, c, d, e, f, g int) {
}
