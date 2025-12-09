// Bad examples for the const001 test case.
package const001

// Bad type for status
type BadStatus int

// Bad: All constants WITHOUT explicit types (violates KTN-CONST-001)
// But respects other rules: proper naming, comments, single-block grouping
const (
	// BAD_MAX_CONNECTIONS defines max connections
	BAD_MAX_CONNECTIONS = 100 // want "KTN-CONST-001"
	// BAD_PORT_NUMBER defines server port
	BAD_PORT_NUMBER = 8080 // want "KTN-CONST-001"
	// BAD_TIMEOUT_MS defines timeout in ms
	BAD_TIMEOUT_MS = 5000 // want "KTN-CONST-001"

	// BAD_HTTP_OK represents HTTP 200 status
	BAD_HTTP_OK = 200 // want "KTN-CONST-001"
	// BAD_HTTP_NOT_FOUND represents 404
	BAD_HTTP_NOT_FOUND = 404 // want "KTN-CONST-001"

	// BAD_API_VERSION defines API version
	BAD_API_VERSION = "v1.0" // want "KTN-CONST-001"
	// BAD_DEFAULT_LANG defines default lang
	BAD_DEFAULT_LANG = "en" // want "KTN-CONST-001"

	// BAD_IS_PRODUCTION indicates prod mode
	BAD_IS_PRODUCTION = true // want "KTN-CONST-001"
	// BAD_ENABLE_CACHE indicates cache on
	BAD_ENABLE_CACHE = false // want "KTN-CONST-001"

	// BAD_RATIO defines calculation ratio
	BAD_RATIO = 1.5 // want "KTN-CONST-001"

	// BAD_STATE_A without explicit type
	BAD_STATE_A = iota // want "KTN-CONST-001"
	// BAD_STATE_B inherits (no error)
	BAD_STATE_B
	// BAD_STATE_C inherits (no error)
	BAD_STATE_C

	// BAD_MULTI_A without explicit type
	// want "KTN-CONST-001"
	// want "KTN-CONST-001"
	BAD_MULTI_A, BAD_MULTI_B = 10, 20
)
