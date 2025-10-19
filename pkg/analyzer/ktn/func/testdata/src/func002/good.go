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

// Good: Example functions are exempt
func ExampleWithManyParams(a, b, c, d, e, f, g int) {
}

// Good: Fuzz functions are exempt
func FuzzWithManyParams(f, a, b, c, d, e, g int) {
}

// Good: Function literal with unnamed parameters (exactly 5)
var GoodLiteralUnnamed = func(int, string, bool, float64, []int) {
}

// Good: Function literal with 4 unnamed parameters
var GoodLiteralFourUnnamed = func(int, int, int, int) {
}

// Good: Function literal with 1 unnamed parameter
var GoodLiteralOneUnnamed = func(int) {
}

// Good: Function literal with no parameters
var GoodLiteralNoParams = func() {
}

// Good: Function with 3 parameters
func ThreeParams(a, b, c int) {
}

// Good: Function with 2 parameters of different types
func TwoParamsMixed(a int, b string) {
}

// Good: Function with 4 parameters grouped by type
func FourParamsGrouped(a, b int, c, d string) {
}

// Good: Function literal with 3 unnamed parameters
var GoodLiteralThreeUnnamed = func(int, string, bool) {
}

// Good: Function literal with 2 unnamed parameters
var GoodLiteralTwoUnnamed = func(int, string) {
}

// Good: Function with variadic parameter counts as 1 param (5 total with variadic)
func WithVariadic(a, b, c, d int, e ...string) {
}
