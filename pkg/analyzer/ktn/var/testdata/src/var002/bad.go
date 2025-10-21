package var002

// Bad: Multiple separate var declarations (violates KTN-VAR-002)
// But respects: VAR-001 (explicit types), VAR-003 (camelCase), VAR-004 (comments), VAR-006 (const before var)

const (
	// BAD_MAX_RETRIES defines maximum retries
	BAD_MAX_RETRIES int = 3
	// PORT_VALUE is port value
	PORT_VALUE int = 8080
	// RATIO_VALUE is ratio value
	RATIO_VALUE float64 = 1.5
)

// badRetries is the first var declaration
var badRetries int = BAD_MAX_RETRIES // want "KTN-VAR-002"

// badConfig is a separate var declaration
var badConfig string = "config" // want "KTN-VAR-002"

// Multiple variables in a group (still separate from other vars)
var ( // want "KTN-VAR-002"
	// badPort is the server port
	badPort int = PORT_VALUE

	// badHost is the server hostname
	badHost string = "localhost"
)

// Another separate var declaration
var badEnabled bool = true // want "KTN-VAR-002"

// Yet another separate var block
var ( // want "KTN-VAR-002"
	// badRatio is a ratio value
	badRatio float64 = RATIO_VALUE

	// badCount is a counter
	badCount int = 0
)
