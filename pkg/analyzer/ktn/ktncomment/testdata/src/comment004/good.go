// Package comment004 provides good test cases.
package comment004

// Good: All variables have comments, explicit types, proper naming, single-block grouping

const (
	// defaultTimeoutValue is the default timeout
	defaultTimeoutValue int = 30
	// httpPortValue is the default HTTP port
	httpPortValue int = 80
	// httpsPortValue is the default HTTPS port
	httpsPortValue int = 443
	// ftpPortValue is the default FTP port
	ftpPortValue int = 21
	// maxConnectionsValue defines maximum concurrent connections
	maxConnectionsValue int = 1000
	// minConnectionsValue defines minimum concurrent connections
	minConnectionsValue int = 10
	// maxRetryCountValue defines maximum number of retries
	maxRetryCountValue int = 5
	// retryDelayMsValue defines delay between retries in milliseconds
	retryDelayMsValue int = 1000
	// backoffMultiplierValue defines the backoff multiplier
	backoffMultiplierValue float64 = 1.5
	// dbPortValue is the database port
	dbPortValue int = 5432
	// onlyConstValue is valid edge case
	onlyConstValue int = 999
)

// All package-level variables grouped in a single block
var (
	// httpPort is the default HTTP port
	httpPort int = httpPortValue

	// httpsPort is the default HTTPS port
	httpsPort int = httpsPortValue

	// ftpPort is the default FTP port
	ftpPort int = ftpPortValue

	// maxConnections defines maximum concurrent connections
	maxConnections int = maxConnectionsValue

	// minConnections defines minimum concurrent connections
	minConnections int = minConnectionsValue

	// defaultTimeout defines the default timeout in seconds
	defaultTimeout int = defaultTimeoutValue

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
	maxRetryCount int = maxRetryCountValue

	// retryDelayMs defines delay between retries in milliseconds
	retryDelayMs int = retryDelayMsValue

	// backoffMultiplier defines the backoff multiplier
	backoffMultiplier float64 = backoffMultiplierValue

	// dbHost is the database host
	dbHost string = "localhost"

	// dbPort is the database port
	dbPort int = dbPortValue

	// dbName is the database name
	dbName string = "mydb"

	// dbUser is the database user
	dbUser string = "admin"
)

// init uses all constants and variables
func init() {
	// Use constants
	_ = onlyConstValue
	// Use variables
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
}
