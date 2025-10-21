package var003

// Good: All package-level variables use camelCase or PascalCase (not SCREAMING_SNAKE_CASE)

const (
	// DEFAULT_TIMEOUT is a constant (SCREAMING_SNAKE_CASE is OK for constants)
	DEFAULT_TIMEOUT int = 30
	// MAX_RETRIES is a constant
	MAX_RETRIES int = 3
	// SERVER_PORT is the server port
	SERVER_PORT int = 8080
	// MAX_CONNECTIONS is maximum connections
	MAX_CONNECTIONS int = 100
	// MIN_TIMEOUT is minimum timeout
	MIN_TIMEOUT int = 5
	// ANSWER is the answer to life
	ANSWER int = 42
)

// All package-level variables grouped in a single block
var (
	// defaultRetries uses camelCase (good for private variables)
	defaultRetries int = DEFAULT_TIMEOUT

	// configuration uses camelCase
	configuration string = "default"

	// isEnabled uses camelCase
	isEnabled bool = false

	// serverPort uses camelCase
	serverPort int = SERVER_PORT

	// maxConnections uses camelCase
	maxConnections int = MAX_CONNECTIONS

	// minTimeout uses camelCase
	minTimeout int = MIN_TIMEOUT

	// apiEndpoint uses camelCase
	apiEndpoint string = "http://localhost"

	// DefaultConfig uses PascalCase (good for exported variables)
	DefaultConfig string = "production"

	// ServerAddress uses PascalCase
	ServerAddress string = "0.0.0.0"
)

// goodFunction demonstrates local variables (not checked by VAR-003).
//
// Returns:
//   - int: calculated value
func goodFunction() int {
	// Local variables are not checked by VAR-003
	localVar := ANSWER
	ANOTHER_LOCAL := MAX_CONNECTIONS
	// Continue traversing AST nodes.
	return localVar + ANOTHER_LOCAL
}
