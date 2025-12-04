// Bad examples for the const001 test case.
package const001

// Bad: All constants in a single grouped block WITHOUT explicit types (violates KTN-CONST-001)
// But respects other rules: proper naming, comments, and single-block grouping
const (
	// MAX_CONNECTIONS defines the maximum number of connections
	MAX_CONNECTIONS = 100
	// PORT_NUMBER defines the server port number
	PORT_NUMBER = 8080
	// TIMEOUT_MS defines the timeout in milliseconds
	TIMEOUT_MS = 5000

	// HTTP_OK represents the HTTP 200 status code
	HTTP_OK = 200
	// HTTP_NOT_FOUND represents the HTTP 404 status code
	HTTP_NOT_FOUND = 404

	// API_VERSION defines the API version string
	API_VERSION = "v1.0"
	// DEFAULT_LANG defines the default language code
	DEFAULT_LANG = "en"

	// IS_PRODUCTION indicates if running in production mode
	IS_PRODUCTION = true
	// ENABLE_CACHE indicates if caching is enabled
	ENABLE_CACHE = false

	// BAD_RATIO defines a calculation ratio (renamed to avoid redeclaration)
	BAD_RATIO = 1.5
)
