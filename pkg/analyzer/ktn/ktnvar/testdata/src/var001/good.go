// Package var001 provides good test cases.
package var001

// Good: All package-level variables use camelCase or PascalCase (not SCREAMING_SNAKE_CASE)

const (
	// DefaultTimeout is a constant in CamelCase
	DefaultTimeout int = 30
	// MaxRetries is a constant
	MaxRetries int = 3
	// ServerPort is the server port
	ServerPort int = 8080
	// MaxConnections is maximum connections
	MaxConnections int = 100
	// MinTimeout is minimum timeout
	MinTimeout int = 5
	// TheAnswer is the answer to life
	TheAnswer int = 42
)

// All package-level variables grouped in a single block
var (
	// defaultRetries uses camelCase (good for private variables)
	defaultRetries int = DefaultTimeout

	// configuration uses camelCase
	configuration string = "default"

	// isEnabled uses camelCase
	isEnabled bool = false

	// serverPort uses camelCase
	serverPort int = ServerPort

	// maxConnections uses camelCase
	maxConnections int = MaxConnections

	// minTimeout uses camelCase
	minTimeout int = MinTimeout

	// apiEndpoint uses camelCase
	apiEndpoint string = "http://localhost"

	// DefaultConfig uses PascalCase (good for exported variables)
	DefaultConfig string = "production"

	// ServerAddress uses PascalCase
	ServerAddress string = "0.0.0.0"

	// ===== Acronyms handling (valid Go naming) =====

	// HTTPClient uses PascalCase with acronym (exported, valid Go convention)
	HTTPClient string = "http-client"

	// httpClient uses camelCase (private, valid Go convention)
	httpClient string = "http-client-private"

	// XMLParser uses PascalCase with acronym (exported)
	XMLParser string = "xml-parser"

	// xmlParser uses camelCase (private)
	xmlParser string = "xml-parser-private"

	// APIEndpoint uses PascalCase with acronym
	APIEndpoint string = "/api/v1"

	// apiKey uses camelCase (private)
	apiKey string = "secret-key"
)

// init utilise les variables pour éviter les erreurs de compilation
func init() {
	// Utilisation des variables privées
	_ = defaultRetries
	_ = configuration
	_ = isEnabled
	_ = serverPort
	_ = maxConnections
	_ = minTimeout
	_ = apiEndpoint
	// Utilisation des variables avec acronymes
	_ = httpClient
	_ = xmlParser
	_ = apiKey
}
