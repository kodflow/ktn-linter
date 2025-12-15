// Package var003 provides good test cases.
package var003

const (
	// MaxRetriesValue is the max retries value
	MaxRetriesValue int = 3
	// TheAnswer is the answer
	TheAnswer int = 42
	// ExampleAge is example age
	ExampleAge int = 30
	// XValue is example x value
	XValue int = 10
	// YValue is example y value
	YValue int = 20
	// Int64Value is int64 example
	Int64Value int64 = 5
	// InitialCount is the initial count
	InitialCount int = 0
)

// Package-level variables can use var with explicit type
var (
	// GlobalConfig is the global configuration
	GlobalConfig string = "default"
	// MaxRetries is the maximum retries
	MaxRetries int = MaxRetriesValue
)

// init demonstrates correct usage of short declarations
func init() {
	// Good: Use := for local variables
	name := "Alice"
	age := ExampleAge
	x, y := XValue, YValue
	_ = name
	_ = age
	_ = x
	_ = y
	// Good: var with explicit type is OK (11 lines)
	var x64 int64 = Int64Value
	var result int
	result = TheAnswer
	_ = x64
	_ = result
	// Good: := for reassign
	count := InitialCount
	count = count + 1
	_ = count
	// Cases where var is NECESSARY (nil values)
	var nilSlice []string
	var nilMap map[string]int
	var nilPointer *int
	var nilChan chan int
	var nilFunc func()
	var nilInterface any
	_ = nilSlice
	_ = nilMap
	_ = nilPointer
	_ = nilChan
	_ = nilFunc
	_ = nilInterface
}
