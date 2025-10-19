package func004

// Bad: Naked return in function with 5 lines
func LongWithNakedReturn() (result int) {
	result = 1
	result += 2
	result += 3
	result += 4
	result += 5
	return // want "KTN-FUNC-004"
}

// Bad: Naked return in longer function
func VeryLongNakedReturn() (a int, b string) {
	a = 1
	a += 2
	a += 3
	a += 4
	a += 5
	a += 6
	b = "test"
	return // want "KTN-FUNC-004"
}

// Bad: Multiple naked returns in long function
func MultipleNakedReturns() (result int) {
	if true {
		result = 1
		result += 2
		result += 3
		result += 4
		result += 5
		return // want "KTN-FUNC-004"
	}
	result = 10
	result += 20
	result += 30
	result += 40
	result += 50
	return // want "KTN-FUNC-004"
}
