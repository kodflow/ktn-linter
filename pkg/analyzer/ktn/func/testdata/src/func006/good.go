package func006

// Good examples: error is always last

func GoodSingleError() error {
	return nil
}

func GoodStringError() (string, error) {
	return "", nil
}

func GoodMultipleReturnsError() (int, string, error) {
	return 0, "", nil
}

func GoodBoolError() (bool, error) {
	return true, nil
}

func GoodNoError() string {
	return ""
}

func GoodNoReturn() {
	// Nothing
}

func GoodMultipleValues() (int, string, bool) {
	return 0, "", false
}

// Method with error last
func (g *GoodType) GoodMethod() (string, error) {
	return "", nil
}

type GoodType struct{}

// Function literal with error last
var goodFunc = func() (int, error) {
	return 0, nil
}

// Closure with error last
func GoodClosure() func() error {
	return func() error {
		return nil
	}
}

// Custom error type (still should be last)
type CustomError struct {
	msg string
}

func (e CustomError) Error() string {
	return e.msg
}

func GoodCustomError() (string, error) {
	return "", CustomError{msg: "test"}
}

// Function with interface{} return (not error)
func GoodInterface() (interface{}, string) {
	return nil, ""
}
