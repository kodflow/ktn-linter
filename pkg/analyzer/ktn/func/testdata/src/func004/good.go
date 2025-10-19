package func004

// Good: No named returns, so naked return doesn't apply
func NoNamedReturns() int {
	return 42
}

// Good: Named returns but very short function (< 5 lines)
func ShortWithNakedReturn() (result int) {
	result = 42
	return
}

// Good: Named returns but returns explicitly
func ExplicitReturn() (result int, err error) {
	result = 42
	err = nil
	return result, err
}

// Good: Named returns in short function (4 lines)
func ShortFourLines() (x int) {
	x = 1
	x = 2
	x = 3
	return
}

// Good: Multiple named returns with explicit return
func MultipleExplicit() (a int, b string, c bool) {
	a = 1
	b = "test"
	c = true
	return a, b, c
}

// Good: Function with no return values at all
func NoReturnValues() {
	x := 1
	_ = x
}

// Good: Function with unnamed return (no named returns, so naked return doesn't apply)
func UnnamedReturn() int {
	return 42
}

// Good: Test function with naked return (exempt)
func TestNakedReturn(t int) (result int) {
	result = 1
	result += 2
	result += 3
	result += 4
	result += 5
	result += 6
	return
}

// Good: Benchmark function with naked return (exempt)
func BenchmarkNakedReturn(b int) (result int) {
	result = 1
	result += 2
	result += 3
	result += 4
	result += 5
	result += 6
	return
}

// Good: Function with returns but all unnamed (so naked return rule doesn't apply)
type Calculator interface {
	Calculate() (int, error)
}

// Good: Function with single unnamed return value
func SingleUnnamedReturn() (int) {
	return 42
}

// Good: Function with multiple unnamed return values
func MultipleUnnamedReturns() (int, string, bool) {
	return 1, "test", true
}
