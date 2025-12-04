// Bad examples for the const004 test case.
package const003

// Bad: Invalid naming (violates KTN-CONST-003)
// But respects: explicit types, comments, and single-block grouping
const (
	// camelCase naming (INVALID)
	// maxSize in camelCase
	maxSize int = 100
	// apiKey in camelCase
	apiKey string = "secret"
	// httpTimeout in camelCase
	httpTimeout int = 30

	// PascalCase naming (INVALID)
	// MaxSize in PascalCase
	MaxSize int = 100
	// ApiKey in PascalCase
	ApiKey string = "secret"
	// HttpTimeout in PascalCase
	HttpTimeout int = 30

	// snake_case (lowercase with underscores) - INVALID
	// max_size in snake_case
	max_size int = 100
	// api_key in snake_case
	api_key string = "secret"
	// http_timeout in snake_case
	http_timeout int = 30

	// Mixed case with underscores - INVALID
	// Max_Size mixed case
	Max_Size int = 100
	// Api_Key mixed case
	Api_Key string = "secret"
	// Http_Timeout mixed case
	Http_Timeout int = 30

	// More camelCase examples
	// statusOk in camelCase
	statusOk int = 200
	// statusCreated in camelCase
	statusCreated int = 201
	// statusError in camelCase
	statusError int = 500

	// More PascalCase examples
	// StateIdle in PascalCase
	StateIdle int = 0
	// StateRunning in PascalCase
	StateRunning int = 1
	// StatePaused in PascalCase
	StatePaused int = 2

	// Mixed variations
	// ErrorNotFound PascalCase
	ErrorNotFound string = "not found"
	// errorUnauthorized camelCase
	errorUnauthorized string = "unauthorized"
	// Error_Internal mixed
	Error_Internal string = "internal"

	// Starting with lowercase
	// defaultPort lowercase start
	defaultPort int = 8080
	// defaultHost lowercase start
	defaultHost string = "localhost"
	// defaultProtocol lowercase start
	defaultProtocol string = "http"

	// Complex camelCase
	// maxConnectionPoolSize complex camelCase
	maxConnectionPoolSize int = 50
	// defaultRequestTimeout complex camelCase
	defaultRequestTimeout int = 60
	// apiKeyHeaderName complex camelCase
	apiKeyHeaderName string = "X-API-Key"

	// Partially correct (mixed) - INVALID
	// MAX_Size partially correct
	MAX_Size int = 100
	// Api_KEY partially correct
	Api_KEY string = "key"
	// HTTP_timeout partially correct
	HTTP_timeout int = 30

	// Database constants with wrong naming
	// dbMaxConnections database setting
	dbMaxConnections int = 100
	// DbMinConnections database setting
	DbMinConnections int = 10
	// db_timeout database setting
	db_timeout int = 30
)
