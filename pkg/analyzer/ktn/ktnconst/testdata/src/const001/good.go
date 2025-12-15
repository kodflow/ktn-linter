// Package const001 provides good test cases.
package const001

// Good: All constants in a single grouped block with explicit types
const (
	// MaxTimeout defines the maximum timeout in seconds
	MaxTimeout int = 30
	// RetryCount specifies the number of retry attempts
	RetryCount int = 3
	// DefaultRatio is the default ratio for calculations
	DefaultRatio float64 = 1.5
	// AppName is the application name
	AppName string = "myapp"

	// Iota with explicit type on first line (using int)
	// StatusPending represents a pending status
	StatusPending int = iota
	// StatusRunning represents a running status (inherits type from previous)
	StatusRunning
	// StatusCompleted represents a completed status (inherits type from previous)
	StatusCompleted

	// Pi represents the mathematical constant pi
	Pi float64 = 3.14159
	// E represents Euler's number
	E float64 = 2.71828

	// FeatureEnabled indicates if the feature is enabled
	FeatureEnabled bool = true
	// DebugMode indicates if debug mode is active
	DebugMode bool = false

	// Exotic types
	// NewlineChar represents a newline character
	NewlineChar rune = '\n'
	// TabChar represents a tab character
	TabChar rune = '\t'
	// ComplexValue represents a complex number
	ComplexValue complex128 = 1 + 2i

	// Multi-name with explicit type (all names share the type)
	// FirstVal and SecondVal are related values
	FirstVal, SecondVal int = 1, 2
)
