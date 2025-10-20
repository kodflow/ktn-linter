package func006

// Bad examples: error is not last

func BadErrorFirst() (error, string) {
	return nil, ""
}

func BadErrorMiddle() (string, error, bool) {
	return "", nil, false
}

func BadErrorFirstOfThree() (error, int, string) {
	return nil, 0, ""
}

// Method with error not last
func (b *BadType) BadMethod() (error, string) {
	return nil, ""
}

type BadType struct{}

// Function literal with error not last
var badFunc = func() (error, int) {
	return nil, 0
}

// Multiple errors with one misplaced (tests the early return after first error found)
func BadMultipleErrors() (error, string, error) {
	return nil, "", nil
}
