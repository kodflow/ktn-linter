package func010

// Bad: 4 unnamed return values
func FourUnnamedReturns() (int, string, bool, error) { // want "KTN-FUNC-010"
	return 1, "test", true, nil
}

// Bad: 5 unnamed return values
func FiveUnnamedReturns() (int, int, string, bool, error) { // want "KTN-FUNC-010"
	return 1, 2, "test", true, nil
}

// Bad: Many unnamed return values
func ManyUnnamedReturns() (int, string, int, bool, float64, error) { // want "KTN-FUNC-010"
	return 1, "test", 30, true, 95.5, nil
}

// Bad: 7 unnamed return values
func SevenUnnamedReturns() (int, string, int, bool, float64, error, string) { // want "KTN-FUNC-010"
	return 1, "test", 30, true, 95.5, nil, "extra"
}
