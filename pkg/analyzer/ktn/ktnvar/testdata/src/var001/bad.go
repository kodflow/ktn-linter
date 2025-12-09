// Bad examples for the var001 test case.
package var001

// Bad: Package-level variables using SCREAMING_SNAKE_CASE (violates KTN-VAR-001)

const (
	// ApiKey is a constant in CamelCase
	ApiKey string = "secret"
	// TimeoutValue is timeout value
	TimeoutValue int = 30
	// PortValue is port value
	PortValue int = 8080
	// MaxConnValue is max connections value
	MaxConnValue int = 100
)

// Variables with SCREAMING_SNAKE_CASE (violates KTN-VAR-001)
var (
	// BAD_TIMEOUT uses SCREAMING_SNAKE_CASE (reserved for constants)
	BAD_TIMEOUT int = TimeoutValue // want "KTN-VAR-001"

	// WRONG_CONFIG uses SCREAMING_SNAKE_CASE (should be camelCase)
	WRONG_CONFIG string = "config" // want "KTN-VAR-001"

	// SERVER_PORT uses SCREAMING_SNAKE_CASE
	SERVER_PORT int = PortValue // want "KTN-VAR-001"

	// SERVER_HOST uses SCREAMING_SNAKE_CASE
	SERVER_HOST string = "localhost" // want "KTN-VAR-001"

	// MAX_CONNECTIONS uses SCREAMING_SNAKE_CASE
	MAX_CONNECTIONS int = MaxConnValue // want "KTN-VAR-001"

	// IS_ENABLED uses SCREAMING_SNAKE_CASE
	IS_ENABLED bool = false // want "KTN-VAR-001"

	// ===== Acronyms with SCREAMING_SNAKE_CASE (should be flagged) =====

	// HTTP_CLIENT uses SCREAMING_SNAKE_CASE with acronym
	HTTP_CLIENT string = "bad" // want "KTN-VAR-001"

	// XML_PARSER uses SCREAMING_SNAKE_CASE with acronym
	XML_PARSER string = "bad" // want "KTN-VAR-001"

	// API_KEY uses SCREAMING_SNAKE_CASE with acronym
	API_KEY string = "bad" // want "KTN-VAR-001"
)
