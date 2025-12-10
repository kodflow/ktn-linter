// Good examples for the var002 test case.
package comment004

// Good: All variables have comments, explicit types, proper naming, single-block grouping

const (
	// DEFAULT_TIMEOUT is the default timeout
	DEFAULT_TIMEOUT int = 30
	// HTTP_PORT is the default HTTP port
	HTTP_PORT int = 80
	// HTTPS_PORT is the default HTTPS port
	HTTPS_PORT int = 443
	// FTP_PORT is the default FTP port
	FTP_PORT int = 21
	// MAX_CONNECTIONS defines maximum concurrent connections
	MAX_CONNECTIONS int = 1000
	// MIN_CONNECTIONS defines minimum concurrent connections
	MIN_CONNECTIONS int = 10
	// MAX_RETRY_COUNT defines maximum number of retries
	MAX_RETRY_COUNT int = 5
	// RETRY_DELAY_MS defines delay between retries in milliseconds
	RETRY_DELAY_MS int = 1000
	// BACKOFF_MULTIPLIER defines the backoff multiplier
	BACKOFF_MULTIPLIER float64 = 1.5
	// DB_PORT is the database port
	DB_PORT int = 5432
	// ONLY_CONST is valid edge case
	ONLY_CONST int = 999
)

// All package-level variables grouped in a single block
var (
	// httpPort is the default HTTP port
	httpPort int = HTTP_PORT

	// httpsPort is the default HTTPS port
	httpsPort int = HTTPS_PORT

	// ftpPort is the default FTP port
	ftpPort int = FTP_PORT

	// maxConnections defines maximum concurrent connections
	maxConnections int = MAX_CONNECTIONS

	// minConnections defines minimum concurrent connections
	minConnections int = MIN_CONNECTIONS

	// defaultTimeout defines the default timeout in seconds
	defaultTimeout int = DEFAULT_TIMEOUT

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
	maxRetryCount int = MAX_RETRY_COUNT

	// retryDelayMs defines delay between retries in milliseconds
	retryDelayMs int = RETRY_DELAY_MS

	// backoffMultiplier defines the backoff multiplier
	backoffMultiplier float64 = BACKOFF_MULTIPLIER

	// dbHost is the database host
	dbHost string = "localhost"

	// dbPort is the database port
	dbPort int = DB_PORT

	// dbName is the database name
	dbName string = "mydb"

	// dbUser is the database user
	dbUser string = "admin"
)
