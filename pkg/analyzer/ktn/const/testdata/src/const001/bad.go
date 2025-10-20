package const001

// Custom type for priority
type Priority int

// Bad: All constants in a single grouped block WITHOUT explicit types (violates KTN-CONST-001)
// But respects other rules: proper naming, comments, and single-block grouping
const (
	// MAX_CONNECTIONS without type
	MAX_CONNECTIONS = 100
	// PORT_NUMBER without type
	PORT_NUMBER = 8080
	// TIMEOUT_MS without type
	TIMEOUT_MS = 5000

	// HTTP_OK status code
	HTTP_OK = 200
	// HTTP_NOT_FOUND status code
	HTTP_NOT_FOUND = 404

	// API_VERSION without type
	API_VERSION = "v1.0"
	// DEFAULT_LANG without type
	DEFAULT_LANG = "en"

	// IS_PRODUCTION flag without type
	IS_PRODUCTION = true
	// ENABLE_CACHE flag without type
	ENABLE_CACHE = false

	// PRIORITY_LOW without explicit type
	PRIORITY_LOW = iota
	// PRIORITY_MEDIUM inherits type (OK)
	PRIORITY_MEDIUM
	// PRIORITY_HIGH inherits type (OK)
	PRIORITY_HIGH
)
