// Package const001 provides good test cases.
package const001

// CustomInt is a custom integer type for testing typed constants.
type CustomInt int

// Good: All constants in a single grouped block with explicit types
const (
	// === Basic types (int, string, bool, float64) ===

	// MaxTimeout defines the maximum timeout in seconds
	MaxTimeout int = 30
	// RetryCount specifies the number of retry attempts
	RetryCount int = 3
	// DefaultRatio is the default ratio for calculations
	DefaultRatio float64 = 1.5
	// AppName is the application name
	AppName string = "myapp"
	// FeatureEnabled indicates if the feature is enabled
	FeatureEnabled bool = true
	// DebugMode indicates if debug mode is active
	DebugMode bool = false

	// === Iota with explicit type ===

	// StatusPending represents a pending status
	StatusPending int = iota
	// StatusRunning represents a running status (inherits type)
	StatusRunning
	// StatusCompleted represents a completed status (inherits type)
	StatusCompleted

	// === Additional numeric types (T1.1) ===

	// ByteValue represents a byte constant
	ByteValue byte = 0xFF
	// Int8Value represents an int8 constant
	Int8Value int8 = -128
	// Int16Value represents an int16 constant
	Int16Value int16 = -32768
	// Int32Value represents an int32 constant
	Int32Value int32 = -2147483648
	// Int64Value represents an int64 constant
	Int64Value int64 = -9223372036854775808
	// UintValue represents a uint constant
	UintValue uint = 42
	// Uint8Value represents a uint8 constant
	Uint8Value uint8 = 255
	// Uint16Value represents a uint16 constant
	Uint16Value uint16 = 65535
	// Uint32Value represents a uint32 constant
	Uint32Value uint32 = 4294967295
	// Uint64Value represents a uint64 constant
	Uint64Value uint64 = 18446744073709551615
	// Float32Value represents a float32 constant
	Float32Value float32 = 3.14
	// UintptrValue represents a uintptr constant
	UintptrValue uintptr = 0x1234

	// === Float and complex types ===

	// Pi represents the mathematical constant pi
	Pi float64 = 3.14159
	// E represents Euler's number
	E float64 = 2.71828
	// Complex64Value represents a complex64 constant
	Complex64Value complex64 = 1 + 2i
	// ComplexValue represents a complex128 number
	ComplexValue complex128 = 1 + 2i

	// === Rune type ===

	// NewlineChar represents a newline character
	NewlineChar rune = '\n'
	// TabChar represents a tab character
	TabChar rune = '\t'

	// === Custom type constant (T1.2) ===

	// CustomValue is a constant with custom type
	CustomValue CustomInt = 42

	// === Multi-name with explicit type ===

	// FirstVal and SecondVal are related values
	FirstVal, SecondVal int = 1, 2
)
