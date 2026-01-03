// Package interface004 contains test cases for KTN-INTERFACE-004.
package interface004

// Processor is a specific interface.
type Processor interface {
	Process() error
}

// ProcessSpecific uses a specific interface.
func ProcessSpecific(p Processor) error {
	return p.Process()
}

// HandleString uses a concrete type.
func HandleString(s string) {
	_ = s
}

// GetInt returns a concrete type.
func GetInt() int {
	return 42
}

// ProcessSlice uses a typed slice.
func ProcessSlice(data []string) {
	_ = data
}
