// Bad examples for the const003 test case.
package const003

// Bad: Invalid naming (violates KTN-CONST-003)
// Contains underscores which is not Go CamelCase convention
const (
	// SCREAMING_SNAKE_CASE naming (INVALID - contains underscores)
	// MaxSizeSnake contains underscore
	MAX_SIZE int = 100 // want "KTN-CONST-003"
	// ApiKeySnake contains underscore
	API_KEY string = "secret" // want "KTN-CONST-003"
	// HttpTimeoutSnake contains underscore
	HTTP_TIMEOUT int = 30 // want "KTN-CONST-003"

	// snake_case (lowercase with underscores) - INVALID
	// maxSizeSnakeLower contains underscore
	max_size int = 100 // want "KTN-CONST-003"
	// apiKeySnakeLower contains underscore
	api_key string = "secret" // want "KTN-CONST-003"
	// httpTimeoutSnakeLower contains underscore
	http_timeout int = 30 // want "KTN-CONST-003"

	// Mixed case with underscores - INVALID
	// MaxSizeMixed contains underscore
	Max_Size int = 100 // want "KTN-CONST-003"
	// ApiKeyMixed contains underscore
	Api_Key string = "secret" // want "KTN-CONST-003"
	// HttpTimeoutMixed contains underscore
	Http_Timeout int = 30 // want "KTN-CONST-003"

	// More SCREAMING_SNAKE examples
	// DatabaseMaxConn contains underscore
	DB_MAX_CONNECTIONS int = 100 // want "KTN-CONST-003"
	// DefaultPortNumber contains underscore
	DEFAULT_PORT int = 8080 // want "KTN-CONST-003"
	// IsProductionMode contains underscore
	IS_PRODUCTION bool = false // want "KTN-CONST-003"

	// Complex underscored names
	// ConnectionPoolSize contains underscore
	MAX_CONNECTION_POOL_SIZE int = 50 // want "KTN-CONST-003"
	// RequestTimeoutSec contains underscore
	DEFAULT_REQUEST_TIMEOUT int = 60 // want "KTN-CONST-003"
	// ApiKeyHeader contains underscore
	API_KEY_HEADER_NAME string = "X-API" // want "KTN-CONST-003"
)
