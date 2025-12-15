// Package var018 provides good test cases.
package var018

// Constants to avoid magic numbers
const (
	// PortValue is the port value
	PortValue int = 8080
	// MaxValue is the max value
	MaxValue int = 100
	// TimeoutValue is the timeout value
	TimeoutValue int = 30
)

// Good: Variables using camelCase or PascalCase (no underscores)

var (
	// httpClient uses camelCase (correct)
	httpClient string = "client"

	// serverPort uses camelCase (correct)
	serverPort int = PortValue

	// maxConnections uses camelCase (correct)
	maxConnections int = MaxValue

	// apiKey uses camelCase (correct)
	apiKey string = "secret"

	// userName uses camelCase (correct)
	userName string = "admin"

	// isEnabled uses camelCase (correct)
	isEnabled bool = true

	// DefaultTimeout uses PascalCase for exported (correct)
	DefaultTimeout int = TimeoutValue

	// HTTPClient uses PascalCase with acronym (correct)
	HTTPClient string = "http"

	// XMLParser uses PascalCase with acronym (correct)
	XMLParser string = "xml"
)

// init uses the variables to avoid compilation errors
func init() {
	_ = httpClient
	_ = serverPort
	_ = maxConnections
	_ = apiKey
	_ = userName
	_ = isEnabled
}
