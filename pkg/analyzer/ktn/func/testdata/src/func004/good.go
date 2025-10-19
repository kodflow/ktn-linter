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
