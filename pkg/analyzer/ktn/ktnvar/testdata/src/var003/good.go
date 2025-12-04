// Good examples for the var003 test case.
package var003

// Good: Variables with correct type visibility

const (
	// DEFAULT_TIMEOUT is the default timeout
	DEFAULT_TIMEOUT int = 30
	// DEFAULT_PORT is the default port
	DEFAULT_PORT int = 8080
	// MAX_CONNECTIONS is the maximum connections
	MAX_CONNECTIONS int = 100
	// ANSWER is the answer to everything
	ANSWER int = 42
	// RETRIES_MULTIPLIER is the multiplier for retries
	RETRIES_MULTIPLIER int = 3
	// BUFFER_SIZE is the buffer size
	BUFFER_SIZE int = 1024
	// BYTE_H is the byte value for H
	BYTE_H byte = 72
	// BYTE_E is the byte value for e
	BYTE_E byte = 101
	// BYTE_L is the byte value for l
	BYTE_L byte = 108
	// BYTE_O is the byte value for o
	BYTE_O byte = 111
	// FLOAT_MULTIPLIER is the float multiplier
	FLOAT_MULTIPLIER int = 3
)

// Cas 1: Type non visible (constante, variable) → type explicite requis
// Cas 2: Type visible dans composite literal → pas de type explicite
// Cas 3: Type visible via conversion de type → pas de type explicite
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
	// endpoints is a list of endpoints
	endpoints = []string{"http://localhost:8080", "http://localhost:9090"}
	// configMap is a map of configuration values
	configMap = map[string]int{"timeout": DEFAULT_TIMEOUT, "retries": RETRIES_MULTIPLIER}
	// buffer is created with make
	buffer = make([]byte, 0, BUFFER_SIZE)
	// cache is a map created with make
	cache = make(map[string]string)
	// convertedInt is a type conversion
	convertedInt = int(ANSWER)
	// convertedStr is a type conversion
	convertedStr = string([]byte{BYTE_H, BYTE_E, BYTE_L, BYTE_L, BYTE_O})
	// convertedFloat is a type conversion
	convertedFloat = float64(FLOAT_MULTIPLIER)
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

// init utilise les fonctions privées
func init() {
	// Appel de goodFunction
	goodFunction()
}
