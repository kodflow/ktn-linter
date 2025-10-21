package func006

// Good examples: error is always last

// GoodSingleError returns only an error (always compliant).
//
// Returns:
//   - error: nil on success
func GoodSingleError() error {
	// Success case
	return nil
}

// GoodStringError returns a string and an error.
//
// Returns:
//   - string: empty string
//   - error: nil on success
func GoodStringError() (string, error) {
	// Success case with empty string
	return "", nil
}

// GoodMultipleReturnsError returns multiple values with error last.
//
// Returns:
//   - int: zero value
//   - string: empty string
//   - error: nil on success
func GoodMultipleReturnsError() (int, string, error) {
	// Success case with zero values
	return 0, "", nil
}

// GoodBoolError returns a boolean and an error.
//
// Returns:
//   - bool: true on success
//   - error: nil on success
func GoodBoolError() (bool, error) {
	// Success case with true
	return true, nil
}

// GoodNoError returns a string without error.
//
// Returns:
//   - string: empty string
func GoodNoError() string {
	// Return empty string
	return ""
}

// GoodNoReturn performs an operation without returning anything.
func GoodNoReturn() {
	// Nothing to do
}

// GoodMultipleValues returns multiple values without error.
//
// Returns:
//   - int: zero value
//   - string: empty string
//   - bool: false
func GoodMultipleValues() (int, string, bool) {
	// Return zero values
	return 0, "", false
}

type GoodType struct{}

// GoodMethod returns a string and an error with error last.
//
// Returns:
//   - string: empty string
//   - error: nil on success
func (g *GoodType) GoodMethod() (string, error) {
	// Success case with empty string
	return "", nil
}

// goodFunc is a function literal that returns an int and an error.
var goodFunc func() (int, error) = func() (int, error) {
	// Success case with zero
	return 0, nil
}

// GoodClosure returns a closure that returns an error.
//
// Returns:
//   - func() error: a function that returns nil error
func GoodClosure() func() error {
	// Return closure that always succeeds
	return func() error {
		// Success case
		return nil
	}
}

// CustomError represents a custom error type.
type CustomError struct {
	msg string
}

// Error returns the error message.
//
// Returns:
//   - string: the error message
func (e CustomError) Error() string {
	// Return stored message
	return e.msg
}

// GoodCustomError returns a string and a custom error.
//
// Returns:
//   - string: empty string
//   - error: custom error with test message
func GoodCustomError() (string, error) {
	// Return custom error
	return "", CustomError{msg: "test"}
}

// GoodInterface returns an interface and a string (not error).
//
// Returns:
//   - interface{}: nil value
//   - string: empty string
func GoodInterface() (interface{}, string) {
	// Return nil interface and empty string
	return nil, ""
}
