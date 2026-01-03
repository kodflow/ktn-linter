// Package interface004 contains test cases for KTN-INTERFACE-004.
package interface004

// ProcessAny uses interface{} as parameter. // want "KTN-INTERFACE-004"
func ProcessAny(data interface{}) {
	_ = data
}

// HandleAny uses any as parameter. // want "KTN-INTERFACE-004"
func HandleAny(input any) {
	_ = input
}

// GetAny returns interface{}. // want "KTN-INTERFACE-004"
func GetAny() interface{} {
	return nil
}

// MultipleEmpty has multiple empty interfaces. // want "KTN-INTERFACE-004" "KTN-INTERFACE-004"
func MultipleEmpty(a interface{}, b any) {
	_, _ = a, b
}
