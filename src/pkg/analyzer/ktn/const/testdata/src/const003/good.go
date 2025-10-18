package const003

// Valid CAPITAL_UNDERSCORE constant names

// Single letter constants
const A = 1
const B = 2
const C = 3

// Acronyms and abbreviations
const API = "api"
const HTTP = "http"
const URL = "url"
const EOF = -1
const HTTPS = "https"

// Properly formatted multi-word constants with underscores
const MAX_SIZE = 100
const API_KEY = "secret"
const HTTP_TIMEOUT = 30
const MAX_BUFFER_SIZE = 1024
const MIN_RETRY_COUNT = 3

// Constants with numbers
const HTTP2 = "http/2"
const TLS1_2 = "tls1.2"
const VERSION_1_0_0 = "1.0.0"
const API_V2_ENDPOINT = "/api/v2"

// Complex multi-word constants
const MAX_CONNECTION_POOL_SIZE = 50
const DEFAULT_REQUEST_TIMEOUT_SECONDS = 60
const API_KEY_HEADER_NAME = "X-API-Key"

// Grouped constants with explicit types
const (
	STATUS_OK      int = 200
	STATUS_CREATED int = 201
	STATUS_ACCEPTED int = 202
)

// Grouped constants with type inference (iota pattern)
const (
	STATE_IDLE int = iota
	STATE_RUNNING
	STATE_PAUSED
	STATE_STOPPED
)

// String constants
const (
	ERROR_MESSAGE_NOT_FOUND     = "not found"
	ERROR_MESSAGE_UNAUTHORIZED  = "unauthorized"
	ERROR_MESSAGE_INTERNAL      = "internal error"
)

// Mixed types in groups
const (
	DEFAULT_PORT     = 8080
	DEFAULT_HOST     = "localhost"
	DEFAULT_PROTOCOL = "http"
	MAX_RETRIES      = 3
)

// Constants with underscores and numbers
const (
	HTTP_1_1_VERSION     = "HTTP/1.1"
	HTTP_2_0_VERSION     = "HTTP/2.0"
	TLS_1_2_CIPHER_SUITE = "TLS_RSA_WITH_AES_128_CBC_SHA"
)

// Database constants
const (
	DB_MAX_CONNECTIONS     = 100
	DB_MIN_CONNECTIONS     = 10
	DB_CONNECTION_TIMEOUT  = 30
	DB_QUERY_TIMEOUT       = 60
)

// File system constants
const (
	FILE_PERMISSIONS      = 0644
	DIR_PERMISSIONS       = 0755
	TEMP_DIR_PREFIX       = "app-temp-"
	LOG_FILE_EXTENSION    = ".log"
)

// Network constants
const (
	TCP_KEEP_ALIVE_INTERVAL = 60
	UDP_BUFFER_SIZE         = 65536
	NETWORK_TIMEOUT_MS      = 5000
)

// Blank identifier (should be skipped)
const _ = "ignored"

// Single character with type
const X int = 10
const Y float64 = 3.14
const Z string = "zed"
