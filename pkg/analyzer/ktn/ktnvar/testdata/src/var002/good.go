// Good examples for the var002 test case.
package var002

// Good: Variables with explicit type AND value (format: var name type = value)

const (
	// DEFAULT_TIMEOUT is the default timeout
	DEFAULT_TIMEOUT int = 30
	// DEFAULT_PORT is the default port
	DEFAULT_PORT int = 8080
	// MAX_CONNECTIONS is the maximum connections
	MAX_CONNECTIONS int = 100
	// ANSWER is the answer to everything
	ANSWER int = 42
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
)

// Style obligatoire: var name type = value
var (
	// defaultRetries has explicit type and value
	defaultRetries int = DEFAULT_TIMEOUT
	// configuration has explicit type and value
	configuration string = "default"
	// isEnabled has explicit type and value
	isEnabled bool = false
	// serverPort has explicit type and value
	serverPort int = DEFAULT_PORT
	// serverHost has explicit type and value
	serverHost string = "localhost"
	// maxConnections has explicit type and value
	maxConnections int = MAX_CONNECTIONS
	// endpoints has explicit type and value
	endpoints []string = []string{"http://localhost:8080", "http://localhost:9090"}
	// configMap has explicit type and value
	configMap map[string]int = map[string]int{"timeout": DEFAULT_TIMEOUT, "retries": 3}
	// buffer has explicit type and value
	buffer []byte = make([]byte, 0, BUFFER_SIZE)
	// cache has explicit type and value
	cache map[string]string = make(map[string]string)
	// convertedInt has explicit type and value
	convertedInt int = int(ANSWER)
	// convertedStr has explicit type and value
	convertedStr string = string([]byte{BYTE_H, BYTE_E, BYTE_L, BYTE_L, BYTE_O})
	// convertedFloat has explicit type and value
	convertedFloat float64 = float64(ANSWER)
)

// goodFunction demonstrates correct local variable usage (not checked by VAR-002).
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
