// Package var013 contains test cases for KTN rules.
package var013

// Bad: Multiple separate var declarations (violates KTN-VAR-013)
// But respects: VAR-001 (explicit types), VAR-003 (camelCase), VAR-004 (comments), VAR-006 (const before var)

const (
	// BadMaxRetries defines maximum retries
	BadMaxRetries int = 3
	// PortValue is port value
	PortValue int = 8080
	// RatioValue is ratio value
	RatioValue float64 = 1.5
)

// badRetries is the first var declaration
var badRetries int = BadMaxRetries

// badConfig is a separate var declaration
var badConfig string = "config"

// Multiple variables in a group (still separate from other vars)
var (
	// badPort is the server port
	badPort int = PortValue

	// badHost is the server hostname
	badHost string = "localhost"
)

// Another separate var declaration
var badEnabled bool = true

// Yet another separate var block
var (
	// badRatio is a ratio value
	badRatio float64 = RatioValue

	// badCount is a counter
	badCount int = 0
)
