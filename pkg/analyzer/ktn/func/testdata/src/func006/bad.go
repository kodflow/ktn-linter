package func006

// Bad examples: error is not last

// BadErrorFirst demonstrates error in first position instead of last.
//
// Returns:
//   - error: potential error (bad position)
//   - string: result message
func BadErrorFirst() (error, string) {
	// Return nil error and empty string
	return nil, ""
}

// BadErrorMiddle demonstrates error in middle position instead of last.
//
// Returns:
//   - string: result message
//   - error: potential error (bad position)
//   - bool: success flag
func BadErrorMiddle() (string, error, bool) {
	// Return empty string, nil error, and false
	return "", nil, false
}

// BadErrorFirstOfThree demonstrates error in first position with three return values.
//
// Returns:
//   - error: potential error (bad position)
//   - int: numeric result
//   - string: result message
func BadErrorFirstOfThree() (error, int, string) {
	// Return nil error, zero, and empty string
	return nil, 0, ""
}

// BadType is a test type for method examples
type BadType struct{}

// BadMethod demonstrates error in first position for a method.
//
// Returns:
//   - error: potential error (bad position)
//   - string: result message
func (b *BadType) BadMethod() (error, string) {
	// Return nil error and empty string
	return nil, ""
}

// badFunc is a function literal with error not last
//
// Returns:
//   - error: potential error (bad position)
//   - int: numeric result
var badFunc func() (error, int) = func() (error, int) {
	// Return nil error and zero
	return nil, 0
}

// BadMultipleErrors demonstrates multiple errors with one misplaced.
//
// Returns:
//   - error: first potential error (bad position)
//   - string: result message
//   - error: second potential error
func BadMultipleErrors() (error, string, error) {
	// Return nil errors and empty string
	return nil, "", nil
}
