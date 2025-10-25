package var001

// Bad: Package-level variables WITHOUT explicit types (violates KTN-VAR-001)
// But respects other rules: proper naming, comments, grouping, after const

const (
	// MAX_RETRIES defines maximum retries
	MAX_RETRIES int = 3
	// PORT_VALUE is port value
	PORT_VALUE int = 8080
	// RATIO_VALUE is ratio value
	RATIO_VALUE float64 = 1.5
)

// Variables without explicit types (violates KTN-VAR-001 only)
var (
	// badRetries has inferred type (should be explicit)
	badRetries = MAX_RETRIES

	// badConfig has inferred type (should be explicit)
	badConfig = "config"

	// badPort has inferred type
	badPort = PORT_VALUE

	// badHost has inferred type
	badHost = "localhost"

	// badEnabled has inferred type
	badEnabled = true

	// badRatio has inferred type
	badRatio = RATIO_VALUE
)
