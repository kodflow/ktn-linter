// Bad examples for the const001 test case.
package const001

// Bad type for status
type BadStatus int

// Bad: All constants WITHOUT explicit types (violates KTN-CONST-001)
// But respects other rules: proper naming, comments, and single-block grouping
const (
	// BAD_MAX_CONNECTIONS defines the maximum number of connections
	BAD_MAX_CONNECTIONS = 100 // want "KTN-CONST-001: la constante 'BAD_MAX_CONNECTIONS' doit avoir un type explicite"
	// BAD_PORT_NUMBER defines the server port number
	BAD_PORT_NUMBER = 8080 // want "KTN-CONST-001: la constante 'BAD_PORT_NUMBER' doit avoir un type explicite"
	// BAD_TIMEOUT_MS defines the timeout in milliseconds
	BAD_TIMEOUT_MS = 5000 // want "KTN-CONST-001: la constante 'BAD_TIMEOUT_MS' doit avoir un type explicite"

	// BAD_HTTP_OK represents the HTTP 200 status code
	BAD_HTTP_OK = 200 // want "KTN-CONST-001: la constante 'BAD_HTTP_OK' doit avoir un type explicite"
	// BAD_HTTP_NOT_FOUND represents the HTTP 404 status code
	BAD_HTTP_NOT_FOUND = 404 // want "KTN-CONST-001: la constante 'BAD_HTTP_NOT_FOUND' doit avoir un type explicite"

	// BAD_API_VERSION defines the API version string
	BAD_API_VERSION = "v1.0" // want "KTN-CONST-001: la constante 'BAD_API_VERSION' doit avoir un type explicite"
	// BAD_DEFAULT_LANG defines the default language code
	BAD_DEFAULT_LANG = "en" // want "KTN-CONST-001: la constante 'BAD_DEFAULT_LANG' doit avoir un type explicite"

	// BAD_IS_PRODUCTION indicates if running in production mode
	BAD_IS_PRODUCTION = true // want "KTN-CONST-001: la constante 'BAD_IS_PRODUCTION' doit avoir un type explicite"
	// BAD_ENABLE_CACHE indicates if caching is enabled
	BAD_ENABLE_CACHE = false // want "KTN-CONST-001: la constante 'BAD_ENABLE_CACHE' doit avoir un type explicite"

	// BAD_RATIO defines a calculation ratio
	BAD_RATIO = 1.5 // want "KTN-CONST-001: la constante 'BAD_RATIO' doit avoir un type explicite"

	// Iota without explicit type (first line must have type)
	// BAD_STATE_A without type
	BAD_STATE_A = iota // want "KTN-CONST-001: la constante 'BAD_STATE_A' doit avoir un type explicite"
	// BAD_STATE_B inherits from previous (no error - no value)
	BAD_STATE_B
	// BAD_STATE_C inherits from previous (no error - no value)
	BAD_STATE_C

	// Multi-name without explicit type
	// BAD_MULTI_A and BAD_MULTI_B without type
	BAD_MULTI_A, BAD_MULTI_B = 10, 20 // want "KTN-CONST-001: la constante 'BAD_MULTI_A' doit avoir un type explicite" "KTN-CONST-001: la constante 'BAD_MULTI_B' doit avoir un type explicite"
)
