package var001

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
)

// Cas 1: Type non visible (constante, variable) → type explicite requis
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

// Cas 2: Type visible dans composite literal → pas de type explicite
var (
	// endpoints is a list of endpoints
	endpoints = []string{"http://localhost:8080", "http://localhost:9090"}
	// configMap is a map of configuration values
	configMap = map[string]int{"timeout": 30, "retries": 3}
	// buffer is created with make
	buffer = make([]byte, 1024)
	// cache is a map created with make
	cache = make(map[string]string)
)

// Cas 3: Type visible via conversion de type → pas de type explicite
var (
	// convertedInt is a type conversion
	convertedInt = int(42)
	// convertedStr is a type conversion
	convertedStr = string([]byte{72, 101, 108, 108, 111})
	// convertedFloat is a type conversion
	convertedFloat = float64(3)
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
