// Good examples for the var014 test case.
package var014

// Good: All package-level variables are grouped in a single var block

const (
	// MAX_RETRIES defines maximum retries
	MAX_RETRIES int = 3
	// DEFAULT_TIMEOUT is the default timeout
	DEFAULT_TIMEOUT int = 30
	// SERVER_PORT is the server port number
	SERVER_PORT int = 8080
	// MAX_CONNECTIONS is the maximum connections
	MAX_CONNECTIONS int = 100
	// ANSWER is the answer to life
	ANSWER int = 42
)

// All variables grouped in a single block
var (
	// defaultRetries defines the default number of retries
	defaultRetries int = MAX_RETRIES

	// configuration holds the app configuration
	configuration string = "default"

	// serverPort is the server port number
	serverPort int = SERVER_PORT

	// serverHost is the server hostname
	serverHost string = "localhost"

	// isEnabled indicates if feature is enabled
	isEnabled bool = false

	// maxConnections is the maximum connections
	maxConnections int = MAX_CONNECTIONS
)

// goodFunction demonstrates local variable usage (not checked by VAR-002).
//
// Returns:
//   - int: calculated value
func goodFunction() int {
	// Local variables are not checked by VAR-002
	localVar := ANSWER
	// Continue traversing AST nodes.
	return localVar
}

// init utilise les fonctions privées
func init() {
	// Appel de goodFunction
	goodFunction()
}
