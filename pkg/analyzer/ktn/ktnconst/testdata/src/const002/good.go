// Package const002 provides good test cases.
package const002

// Good: All constants grouped in a single block at the very top
// Order is correct: const → var → type
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
)

// Variables come after constants (correct order)
var (
	// GlobalVar1 is the first global variable
	GlobalVar1 string = "var1"
	// GlobalVar2 is the second global variable
	GlobalVar2 string = "var2"
)

// Types come after variables (correct order)
type goodType struct {
	// Field is a field
	Field string
}

// anotherType is another type declaration
type anotherType int
