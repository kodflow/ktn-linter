package const004

// Good: All constants have comments, explicit types, proper naming, single-block grouping
const (
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
	// DEFAULT_TIMEOUT defines the default timeout in seconds
	DEFAULT_TIMEOUT int = 30

	// API_VERSION is the current API version
	API_VERSION string = "v1.0"
	// API_ENDPOINT is the base API endpoint
	API_ENDPOINT string = "/api"
	// API_KEY is the authentication key
	API_KEY string = "secret"

	// FEATURE_ENABLED indicates if the feature is enabled
	FEATURE_ENABLED bool = true
	// DEBUG_MODE indicates if debug mode is active
	DEBUG_MODE bool = false
	// VERBOSE_LOGGING indicates if verbose logging is enabled
	VERBOSE_LOGGING bool = false

	// MAX_RETRY_COUNT defines maximum number of retries
	MAX_RETRY_COUNT int = 5
	// RETRY_DELAY_MS defines delay between retries in milliseconds
	RETRY_DELAY_MS int = 1000
	// BACKOFF_MULTIPLIER defines the backoff multiplier
	BACKOFF_MULTIPLIER float64 = 1.5

	// DB_HOST is the database host
	DB_HOST string = "localhost"
	// DB_PORT is the database port
	DB_PORT int = 5432
	// DB_NAME is the database name
	DB_NAME string = "mydb"
	// DB_USER is the database user
	DB_USER string = "admin"
)
