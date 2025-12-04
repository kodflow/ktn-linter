// Good examples for the var001 test case.
package var001

// Good: All package-level variables use camelCase or PascalCase (not SCREAMING_SNAKE_CASE)

const (
	// GOOD_DEFAULT_TIMEOUT is a constant (SCREAMING_SNAKE_CASE is OK for constants)
	GOOD_DEFAULT_TIMEOUT int = 30
	// GOOD_MAX_RETRIES is a constant
	GOOD_MAX_RETRIES int = 3
	// GOOD_SERVER_PORT is the server port
	GOOD_SERVER_PORT int = 8080
	// GOOD_MAX_CONNECTIONS is maximum connections
	GOOD_MAX_CONNECTIONS int = 100
	// GOOD_MIN_TIMEOUT is minimum timeout
	GOOD_MIN_TIMEOUT int = 5
	// GOOD_ANSWER is the answer to life
	GOOD_ANSWER int = 42
)

// All package-level variables grouped in a single block
var (
	// defaultRetries uses camelCase (good for private variables)
	defaultRetries int = GOOD_DEFAULT_TIMEOUT

	// configuration uses camelCase
	configuration string = "default"

	// isEnabled uses camelCase
	isEnabled bool = false

	// serverPort uses camelCase
	serverPort int = GOOD_SERVER_PORT

	// maxConnections uses camelCase
	maxConnections int = GOOD_MAX_CONNECTIONS

	// minTimeout uses camelCase
	minTimeout int = GOOD_MIN_TIMEOUT

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
	localVar := GOOD_ANSWER
	ANOTHER_LOCAL := GOOD_MAX_CONNECTIONS
	// Continue traversing AST nodes.
	return localVar + ANOTHER_LOCAL
}

// init utilise les fonctions privées
func init() {
	// Appel de goodFunction
	goodFunction()
}
