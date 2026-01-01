// Package comment003 provides good test cases.
package comment003

// Good: All constants have comments, explicit types, proper naming, single-block grouping
const (
	// httpPort is the default HTTP port
	httpPort int = 80
	// httpsPort is the default HTTPS port
	httpsPort int = 443
	// ftpPort is the default FTP port
	ftpPort int = 21

	// maxConnections defines maximum concurrent connections
	maxConnections int = 1000
	// minConnections defines minimum concurrent connections
	minConnections int = 10
	// defaultTimeout defines the default timeout in seconds
	defaultTimeout int = 30

	// apiVersion is the current API version
	apiVersion string = "v1.0"
	// apiEndpoint is the base API endpoint
	apiEndpoint string = "/api"
	// apiKey is the authentication key
	apiKey string = "secret"

	// featureEnabled indicates if the feature is enabled
	featureEnabled bool = true
	// debugMode indicates if debug mode is active
	debugMode bool = false
	// verboseLogging indicates if verbose logging is enabled
	verboseLogging bool = false

	// maxRetryCount defines maximum number of retries
	maxRetryCount int = 5
	// retryDelayMs defines delay between retries in milliseconds
	retryDelayMs int = 1000
	// backoffMultiplier defines the backoff multiplier
	backoffMultiplier float64 = 1.5

	// dbHost is the database host
	dbHost string = "localhost"
	// dbPort is the database port
	dbPort int = 5432
	// dbName is the database name
	dbName string = "mydb"
	// dbUser is the database user
	dbUser string = "admin"
)

// onlyVar is a test variable
var onlyVar string = "test"

// init uses the constants to prevent unused errors
func init() {
	// Use all constants
	_ = httpPort
	_ = httpsPort
	_ = ftpPort
	_ = maxConnections
	_ = minConnections
	_ = defaultTimeout
	_ = apiVersion
	_ = apiEndpoint
	_ = apiKey
	_ = featureEnabled
	_ = debugMode
	_ = verboseLogging
	_ = maxRetryCount
	_ = retryDelayMs
	_ = backoffMultiplier
	_ = dbHost
	_ = dbPort
	_ = dbName
	_ = dbUser
	_ = onlyVar
}
