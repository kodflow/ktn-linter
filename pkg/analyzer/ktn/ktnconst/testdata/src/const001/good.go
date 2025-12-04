// Good examples for the const001 test case.
package const001

// Custom type for status
type Status int

// Good: All constants in a single grouped block with explicit types
const (
	// MAX_TIMEOUT defines the maximum timeout in seconds
	MAX_TIMEOUT int = 30
	// RETRY_COUNT specifies the number of retry attempts
	RETRY_COUNT int = 3
	// DEFAULT_RATIO is the default ratio for calculations
	DEFAULT_RATIO float64 = 1.5
	// APP_NAME is the application name
	APP_NAME string = "myapp"

	// Iota with explicit type on first line
	// STATUS_PENDING represents a pending status
	STATUS_PENDING Status = iota
	// STATUS_RUNNING represents a running status (inherits type from previous)
	STATUS_RUNNING
	// STATUS_COMPLETED represents a completed status (inherits type from previous)
	STATUS_COMPLETED

	// PI represents the mathematical constant pi
	PI float64 = 3.14159
	// E represents Euler's number
	E float64 = 2.71828

	// FEATURE_ENABLED indicates if the feature is enabled
	FEATURE_ENABLED bool = true
	// DEBUG_MODE indicates if debug mode is active
	DEBUG_MODE bool = false

	// Exotic types
	// NEWLINE_CHAR represents a newline character
	NEWLINE_CHAR rune = '\n'
	// TAB_CHAR represents a tab character
	TAB_CHAR rune = '\t'
	// COMPLEX_VALUE represents a complex number
	COMPLEX_VALUE complex128 = 1 + 2i

	// Multi-name with explicit type (all names share the type)
	// FIRST_VAL and SECOND_VAL are related values
	FIRST_VAL, SECOND_VAL int = 1, 2
)
