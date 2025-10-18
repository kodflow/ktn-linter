package const003

// Invalid constant names - not following CAPITAL_UNDERSCORE convention

// camelCase - INVALID
const maxSize = 100              // want `KTN-CONST-003: la constante 'maxSize' doit utiliser la convention CAPITAL_UNDERSCORE`
const apiKey = "secret"          // want `KTN-CONST-003: la constante 'apiKey' doit utiliser la convention CAPITAL_UNDERSCORE`
const httpTimeout = 30           // want `KTN-CONST-003: la constante 'httpTimeout' doit utiliser la convention CAPITAL_UNDERSCORE`

// PascalCase - INVALID
const MaxSize = 100              // want `KTN-CONST-003: la constante 'MaxSize' doit utiliser la convention CAPITAL_UNDERSCORE`
const ApiKey = "secret"          // want `KTN-CONST-003: la constante 'ApiKey' doit utiliser la convention CAPITAL_UNDERSCORE`
const HttpTimeout = 30           // want `KTN-CONST-003: la constante 'HttpTimeout' doit utiliser la convention CAPITAL_UNDERSCORE`

// snake_case (lowercase with underscores) - INVALID
const max_size = 100             // want `KTN-CONST-003: la constante 'max_size' doit utiliser la convention CAPITAL_UNDERSCORE`
const api_key = "secret"         // want `KTN-CONST-003: la constante 'api_key' doit utiliser la convention CAPITAL_UNDERSCORE`
const http_timeout = 30          // want `KTN-CONST-003: la constante 'http_timeout' doit utiliser la convention CAPITAL_UNDERSCORE`

// Mixed case with underscores - INVALID
const Max_Size = 100             // want `KTN-CONST-003: la constante 'Max_Size' doit utiliser la convention CAPITAL_UNDERSCORE`
const Api_Key = "secret"         // want `KTN-CONST-003: la constante 'Api_Key' doit utiliser la convention CAPITAL_UNDERSCORE`
const Http_Timeout = 30          // want `KTN-CONST-003: la constante 'Http_Timeout' doit utiliser la convention CAPITAL_UNDERSCORE`

// Grouped constants with invalid names
const (
	statusOk      = 200           // want `KTN-CONST-003: la constante 'statusOk' doit utiliser la convention CAPITAL_UNDERSCORE`
	statusCreated = 201           // want `KTN-CONST-003: la constante 'statusCreated' doit utiliser la convention CAPITAL_UNDERSCORE`
	statusError   = 500           // want `KTN-CONST-003: la constante 'statusError' doit utiliser la convention CAPITAL_UNDERSCORE`
)

// PascalCase in groups - INVALID
const (
	StateIdle    = 0              // want `KTN-CONST-003: la constante 'StateIdle' doit utiliser la convention CAPITAL_UNDERSCORE`
	StateRunning = 1              // want `KTN-CONST-003: la constante 'StateRunning' doit utiliser la convention CAPITAL_UNDERSCORE`
	StatePaused  = 2              // want `KTN-CONST-003: la constante 'StatePaused' doit utiliser la convention CAPITAL_UNDERSCORE`
)

// Mixed case variations - INVALID
const (
	ErrorNotFound     = "not found"      // want `KTN-CONST-003: la constante 'ErrorNotFound' doit utiliser la convention CAPITAL_UNDERSCORE`
	errorUnauthorized = "unauthorized"   // want `KTN-CONST-003: la constante 'errorUnauthorized' doit utiliser la convention CAPITAL_UNDERSCORE`
	Error_Internal    = "internal"       // want `KTN-CONST-003: la constante 'Error_Internal' doit utiliser la convention CAPITAL_UNDERSCORE`
)

// Starting with lowercase - INVALID
const (
	defaultPort     = 8080        // want `KTN-CONST-003: la constante 'defaultPort' doit utiliser la convention CAPITAL_UNDERSCORE`
	defaultHost     = "localhost" // want `KTN-CONST-003: la constante 'defaultHost' doit utiliser la convention CAPITAL_UNDERSCORE`
	defaultProtocol = "http"      // want `KTN-CONST-003: la constante 'defaultProtocol' doit utiliser la convention CAPITAL_UNDERSCORE`
)

// Mixed naming in single group - INVALID
const (
	validName        = "OK"        // want `KTN-CONST-003: la constante 'validName' doit utiliser la convention CAPITAL_UNDERSCORE`
	AnotherBadName   = "BAD"       // want `KTN-CONST-003: la constante 'AnotherBadName' doit utiliser la convention CAPITAL_UNDERSCORE`
	yet_another_bad  = "WORSE"     // want `KTN-CONST-003: la constante 'yet_another_bad' doit utiliser la convention CAPITAL_UNDERSCORE`
)

// Complex camelCase - INVALID
const (
	maxConnectionPoolSize = 50             // want `KTN-CONST-003: la constante 'maxConnectionPoolSize' doit utiliser la convention CAPITAL_UNDERSCORE`
	defaultRequestTimeout = 60             // want `KTN-CONST-003: la constante 'defaultRequestTimeout' doit utiliser la convention CAPITAL_UNDERSCORE`
	apiKeyHeaderName      = "X-API-Key"    // want `KTN-CONST-003: la constante 'apiKeyHeaderName' doit utiliser la convention CAPITAL_UNDERSCORE`
)

// Complex PascalCase - INVALID
const (
	MaxConnectionPoolSize = 50             // want `KTN-CONST-003: la constante 'MaxConnectionPoolSize' doit utiliser la convention CAPITAL_UNDERSCORE`
	DefaultRequestTimeout = 60             // want `KTN-CONST-003: la constante 'DefaultRequestTimeout' doit utiliser la convention CAPITAL_UNDERSCORE`
	ApiKeyHeaderName      = "X-API-Key"    // want `KTN-CONST-003: la constante 'ApiKeyHeaderName' doit utiliser la convention CAPITAL_UNDERSCORE`
)

// Numbers with wrong case - INVALID
const (
	http2Version = "HTTP/2.0"      // want `KTN-CONST-003: la constante 'http2Version' doit utiliser la convention CAPITAL_UNDERSCORE`
	tls1_2       = "TLS 1.2"       // want `KTN-CONST-003: la constante 'tls1_2' doit utiliser la convention CAPITAL_UNDERSCORE`
	apiV2Url     = "/api/v2"       // want `KTN-CONST-003: la constante 'apiV2Url' doit utiliser la convention CAPITAL_UNDERSCORE`
)

// Partially correct (mixed uppercase/lowercase) - INVALID
const (
	MAX_Size        = 100          // want `KTN-CONST-003: la constante 'MAX_Size' doit utiliser la convention CAPITAL_UNDERSCORE`
	Api_KEY         = "key"        // want `KTN-CONST-003: la constante 'Api_KEY' doit utiliser la convention CAPITAL_UNDERSCORE`
	HTTP_timeout    = 30           // want `KTN-CONST-003: la constante 'HTTP_timeout' doit utiliser la convention CAPITAL_UNDERSCORE`
)

// Database constants with wrong naming - INVALID
const (
	dbMaxConnections = 100         // want `KTN-CONST-003: la constante 'dbMaxConnections' doit utiliser la convention CAPITAL_UNDERSCORE`
	DbMinConnections = 10          // want `KTN-CONST-003: la constante 'DbMinConnections' doit utiliser la convention CAPITAL_UNDERSCORE`
	db_timeout       = 30          // want `KTN-CONST-003: la constante 'db_timeout' doit utiliser la convention CAPITAL_UNDERSCORE`
)
