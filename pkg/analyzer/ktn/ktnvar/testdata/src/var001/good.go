package var001

// Good: All package-level variables have explicit types

const (
	// DEFAULT_TIMEOUT is the default timeout
	DEFAULT_TIMEOUT int = 30
	// DEFAULT_PORT is the default port
	DEFAULT_PORT int = 8080
	// MAX_CONNECTIONS is the maximum connections
	MAX_CONNECTIONS int = 100
	// ANSWER is the answer to everything
	ANSWER int = 42
)

// All variables grouped in a single block
var (
	// defaultRetries defines the default number of retries
	defaultRetries int = DEFAULT_TIMEOUT
	// configuration holds the app configuration
	configuration string = "default"
	// isEnabled indicates if feature is enabled
	isEnabled bool = false
	// serverPort is the server port number
	serverPort int = DEFAULT_PORT
	// serverHost is the server hostname
	serverHost string = "localhost"
	// maxConnections is the maximum connections
	maxConnections int = MAX_CONNECTIONS
)

// goodFunction demonstrates correct local variable usage (not checked by VAR-001).
//
// Returns:
//   - int: calculated value
func goodFunction() int {
	// Local variables are not checked by VAR-001
	localVar := ANSWER
	// Continue traversing AST nodes.
	return localVar
}
