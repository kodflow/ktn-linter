// Package shared provides common utilities for static analysis.
package shared

// TestTarget represents the target of a test function.
// It contains parsed information from a test function name to determine
// which function or method the test is intended to cover.
type TestTarget struct {
	// FuncName is the name of the function being tested.
	FuncName string
	// ReceiverName is the receiver type name (for method tests).
	ReceiverName string
	// IsPrivate indicates if this targets a private function.
	IsPrivate bool
	// IsMethod indicates if this targets a method.
	IsMethod bool
}
