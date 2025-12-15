// Package var013 provides good test cases.
package var013

const (
	// MaxRetries defines maximum retries
	MaxRetries int = 3
	// DefaultTimeout is the default timeout
	DefaultTimeout int = 30
	// ServerPort is the server port number
	ServerPort int = 8080
	// MaxConnections is the maximum connections
	MaxConnections int = 100
)

// All variables grouped in a single block
var (
	// defaultRetries defines the default number of retries
	defaultRetries int = MaxRetries

	// configuration holds the app configuration
	configuration string = "default"

	// serverPort is the server port number
	serverPort int = ServerPort

	// serverHost is the server hostname
	serverHost string = "localhost"

	// isEnabled indicates if feature is enabled
	isEnabled bool = false

	// maxConnections is the maximum connections
	maxConnections int = MaxConnections
)

// init demonstrates correct usage patterns
func init() {
	// Local variables are not checked by VAR-013
	localVar := MaxRetries
	_ = localVar
	_ = defaultRetries
	_ = configuration
	_ = serverPort
	_ = serverHost
	_ = isEnabled
	_ = maxConnections
}
