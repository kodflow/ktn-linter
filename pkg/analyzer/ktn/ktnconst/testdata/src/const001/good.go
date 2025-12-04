// Good examples for the const001 test case.
package const001

// Custom type for status
type Status int

// Good: All constants in a single grouped block with explicit types, proper naming, and comments
const (
	// MAX_TIMEOUT defines the maximum timeout in seconds
	MAX_TIMEOUT int = 30
	// RETRY_COUNT specifies the number of retry attempts
	RETRY_COUNT int = 3
	// DEFAULT_RATIO is the default ratio for calculations
	DEFAULT_RATIO float64 = 1.5
	// APP_NAME is the application name
	APP_NAME string = "myapp"

	// STATUS_PENDING represents a pending status
	STATUS_PENDING Status = iota
	// STATUS_RUNNING represents a running status
	STATUS_RUNNING
	// STATUS_COMPLETED represents a completed status
	STATUS_COMPLETED

	// PI represents the mathematical constant π
	PI float64 = 3.14159
	// E represents Euler's number
	E float64 = 2.71828

	// FEATURE_ENABLED indicates if the feature is enabled
	FEATURE_ENABLED bool = true
	// DEBUG_MODE indicates if debug mode is active
	DEBUG_MODE bool = false
)
