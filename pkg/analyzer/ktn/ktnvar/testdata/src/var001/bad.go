// Bad examples for the var001 test case.
package var001

// Bad: Package-level variables using SCREAMING_SNAKE_CASE (violates KTN-VAR-001)
// But respects other rules: explicit types (VAR-001), proper comments (VAR-004), grouped (VAR-002)

const (
	// API_KEY is a constant
	API_KEY string = "secret"
	// TIMEOUT_VALUE is timeout value
	TIMEOUT_VALUE int = 30
	// PORT_VALUE is port value
	PORT_VALUE int = 8080
	// MAX_CONN_VALUE is max connections value
	MAX_CONN_VALUE int = 100
)

// Variables with SCREAMING_SNAKE_CASE (violates KTN-VAR-001)
var (
	// BAD_TIMEOUT uses SCREAMING_SNAKE_CASE (reserved for constants)
	BAD_TIMEOUT int = TIMEOUT_VALUE

	// WRONG_CONFIG uses SCREAMING_SNAKE_CASE (should be camelCase)
	WRONG_CONFIG string = "config"

	// SERVER_PORT uses SCREAMING_SNAKE_CASE
	SERVER_PORT int = PORT_VALUE

	// SERVER_HOST uses SCREAMING_SNAKE_CASE
	SERVER_HOST string = "localhost"

	// MAX_CONNECTIONS uses SCREAMING_SNAKE_CASE
	MAX_CONNECTIONS int = MAX_CONN_VALUE

	// IS_ENABLED uses SCREAMING_SNAKE_CASE
	IS_ENABLED bool = false
)
