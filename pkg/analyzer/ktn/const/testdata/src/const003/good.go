package const003

// Good: All constants use CAPITAL_UNDERSCORE naming, explicit types, comments, single block
const (
	// Single letter constants are valid
	// A represents the first value
	A int = 1
	// B represents the second value
	B int = 2
	// C represents the third value
	C int = 3

	// Acronyms are valid
	// API endpoint path
	API string = "api"
	// HTTP protocol
	HTTP string = "http"
	// URL format
	URL string = "url"
	// EOF end of file marker
	EOF int = -1
	// HTTPS secure protocol
	HTTPS string = "https"

	// Multi-word constants with underscores
	// MAX_SIZE defines the maximum size
	MAX_SIZE int = 100
	// API_KEY is the authentication key
	API_KEY string = "secret"
	// HTTP_TIMEOUT defines timeout in seconds
	HTTP_TIMEOUT int = 30
	// MAX_BUFFER_SIZE defines buffer size
	MAX_BUFFER_SIZE int = 1024
	// MIN_RETRY_COUNT defines minimum retries
	MIN_RETRY_COUNT int = 3

	// Constants with numbers
	// HTTP2 protocol version
	HTTP2 string = "http/2"
	// TLS1_2 TLS version
	TLS1_2 string = "tls1.2"
	// VERSION_1_0_0 software version
	VERSION_1_0_0 string = "1.0.0"
	// API_V2_ENDPOINT version 2 endpoint
	API_V2_ENDPOINT string = "/api/v2"

	// Complex multi-word constants
	// MAX_CONNECTION_POOL_SIZE defines pool size
	MAX_CONNECTION_POOL_SIZE int = 50
	// DEFAULT_REQUEST_TIMEOUT_SECONDS timeout value
	DEFAULT_REQUEST_TIMEOUT_SECONDS int = 60
	// API_KEY_HEADER_NAME header name for API key
	API_KEY_HEADER_NAME string = "X-API-Key"

	// Status codes
	// STATUS_OK success status
	STATUS_OK int = 200
	// STATUS_CREATED resource created
	STATUS_CREATED int = 201
	// STATUS_ACCEPTED request accepted
	STATUS_ACCEPTED int = 202

	// State constants
	// STATE_IDLE idle state
	STATE_IDLE int = 0
	// STATE_RUNNING running state
	STATE_RUNNING int = 1
	// STATE_PAUSED paused state
	STATE_PAUSED int = 2

	// Database settings
	// DB_MAX_CONNECTIONS maximum connections
	DB_MAX_CONNECTIONS int = 100
	// DB_MIN_CONNECTIONS minimum connections
	DB_MIN_CONNECTIONS int = 10
	// DB_TIMEOUT database timeout
	DB_TIMEOUT int = 30

	// Blank identifier (valid edge case - should be skipped by naming rules)
	_ int = 999
)
