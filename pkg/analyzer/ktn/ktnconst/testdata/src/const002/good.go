// Package const002 provides good test cases.
package const002

// Good: All constants grouped in a single block at the very top
// Order is correct: const → var → type → func
const (
	// ConfigValue1 is the first configuration value
	ConfigValue1 string = "config1"
	// ConfigValue2 is the second configuration value
	ConfigValue2 string = "config2"
	// ConfigValue3 is the third configuration value
	ConfigValue3 string = "config3"
	// MaxRetry defines the maximum retry count
	MaxRetry int = 5
	// TimeoutSec defines the timeout in seconds
	TimeoutSec int = 30

	// === Iota pattern with int type (T2.2) ===
	// Demonstrates iota within the main const block

	// StatusPending is the pending status
	StatusPending int = iota
	// StatusActive is the active status
	StatusActive
	// StatusDone is the done status
	StatusDone
)

// Variables come after constants (correct order)
var (
	// GlobalVar1 is the first global variable
	GlobalVar1 string = "var1"
	// GlobalVar2 is the second global variable
	GlobalVar2 string = "var2"
)

// === Type declarations come after variables (correct order) ===

// goodType is a struct type
type goodType struct {
	// Field is a field
	Field string
}

// anotherType is another type declaration
type anotherType int

// Status is a custom type for status values
type Status int

// === Functions come after types (correct order - T2.1) ===

// init demonstrates func declarations after types (correct order).
//
// Params: none
//
// Returns: none
func init() {
	// Use declarations to avoid unused errors
	_ = ConfigValue1
	_ = GlobalVar1
	_ = goodType{}
	_ = anotherType(0)
	_ = Status(0)
}
