package func007

// Missing documentation completely
func NoDoc() string { // want "KTN-FUNC-007"
	return ""
}

// Wrong format - doesn't start with function name
// This is wrong
func WrongFormat() { // want "KTN-FUNC-007"
	// Body
}

// BadMissingParams has params but no Params section
func BadMissingParams(x int) { // want "KTN-FUNC-007"
	// Body
}

// BadMissingReturns has returns but no Returns section
func BadMissingReturns() string { // want "KTN-FUNC-007"
	return ""
}
