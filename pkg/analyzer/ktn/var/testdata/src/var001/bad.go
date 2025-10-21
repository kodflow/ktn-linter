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
	badRetries = MAX_RETRIES // want "KTN-VAR-001"

	// badConfig has inferred type (should be explicit)
	badConfig = "config" // want "KTN-VAR-001"

	// badPort has inferred type
	badPort = PORT_VALUE // want "KTN-VAR-001"

	// badHost has inferred type
	badHost = "localhost" // want "KTN-VAR-001"

	// badEnabled has inferred type
	badEnabled = true // want "KTN-VAR-001"

	// badRatio has inferred type
	badRatio = RATIO_VALUE // want "KTN-VAR-001"
)
