package func006

// Bad examples: error is not last

func BadErrorFirst() (error, string) { // want "KTN-FUNC-006"
	return nil, ""
}

func BadErrorMiddle() (string, error, bool) { // want "KTN-FUNC-006"
	return "", nil, false
}

func BadErrorFirstOfThree() (error, int, string) { // want "KTN-FUNC-006"
	return nil, 0, ""
}

// Method with error not last
func (b *BadType) BadMethod() (error, string) { // want "KTN-FUNC-006"
	return nil, ""
}

type BadType struct{}

// Function literal with error not last
var badFunc = func() (error, int) { // want "KTN-FUNC-006"
	return nil, 0
}
