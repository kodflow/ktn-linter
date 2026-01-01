// Package const003 contains test cases for KTN-CONST-003.
// This file contains ALL cases that MUST NOT trigger KTN-CONST-003 errors.
// Go convention: Use CamelCase (MixedCaps), no underscores.
package const003

// =============================================================================
// SECTION 1: Single letter constants (exported and unexported)
// =============================================================================

const (
	// A represents the first value (exported)
	A int = 1
	// B represents the second value (exported)
	B int = 2
	// C represents the third value (exported)
	C int = 3
	// X coordinate
	X int = 10
	// Y coordinate
	Y int = 20
	// Z coordinate
	Z int = 30
)

const (
	// i is a common loop variable
	i int = 0
	// j is a common loop variable
	j int = 1
	// k is a common loop variable
	k int = 2
	// n is a common count variable
	n int = 100
	// x is a common variable
	x int = 42
	// y is a common variable
	y int = 43
)

// =============================================================================
// SECTION 2: PascalCase for exported constants
// =============================================================================

const (
	// MaxSize defines the maximum size
	MaxSize int = 100
	// ApiKey is the authentication key
	ApiKey string = "secret"
	// HttpTimeout defines timeout in seconds
	HttpTimeout int = 30
	// MaxBufferSize defines buffer size
	MaxBufferSize int = 1024
	// MinRetryCount defines minimum retries
	MinRetryCount int = 3
	// DefaultPort default server port
	DefaultPort int = 8080
	// IsProduction production mode flag
	IsProduction bool = false
)

// =============================================================================
// SECTION 3: camelCase for unexported constants
// =============================================================================

const (
	// maxInternalSize is internal max size
	maxInternalSize int = 50
	// defaultTimeout is the default timeout
	defaultTimeout int = 10
	// httpInternalPort is internal port
	httpInternalPort int = 8080
	// apiSecretKey is the secret API key
	apiSecretKey string = "secret123"
	// isDebugMode debug mode flag
	isDebugMode bool = true
	// connectionPoolSize pool size
	connectionPoolSize int = 25
)

// =============================================================================
// SECTION 4: Acronyms in correct Go style (all caps as one unit)
// =============================================================================

const (
	// APIEndpoint is the API endpoint
	APIEndpoint string = "/api"
	// HTTPStatus is the HTTP status code
	HTTPStatus int = 200
	// URLPath is the URL path
	URLPath string = "/path"
	// XMLParser is the XML parser name
	XMLParser string = "xml"
	// JSONFormat is JSON format
	JSONFormat string = "json"
	// HTMLTemplate is HTML template
	HTMLTemplate string = "html"
	// TCPProtocol is TCP protocol
	TCPProtocol string = "tcp"
	// UDPProtocol is UDP protocol
	UDPProtocol string = "udp"
)

// =============================================================================
// SECTION 5: Acronyms at different positions
// =============================================================================

const (
	// appID is application ID (ID at end)
	appID string = "app123"
	// userID is user ID
	userID string = "user456"
	// xmlHTTPRequest is XMLHttpRequest pattern
	xmlHTTPRequest string = "xhr"
	// httpAPIClient is HTTP API client
	httpAPIClient string = "client"
	// ServeHTTP is serve HTTP pattern
	ServeHTTP string = "serve"
	// ParseURL is parse URL pattern
	ParseURL string = "parse"
)

// =============================================================================
// SECTION 6: Constants with numbers
// =============================================================================

const (
	// Http2Protocol is HTTP/2 protocol
	Http2Protocol string = "h2"
	// Tls12Version is TLS 1.2 version
	Tls12Version string = "1.2"
	// Version100 is version 1.0.0
	Version100 string = "1.0.0"
	// Config2 second config
	Config2 int = 2
	// V3 version 3
	V3 int = 3
	// A1B2C3 mixed letters and numbers
	A1B2C3 int = 123
	// Port8080 port number
	Port8080 int = 8080
	// Error404 HTTP 404 error
	Error404 int = 404
)

// =============================================================================
// SECTION 7: Complex names without underscores
// =============================================================================

const (
	// MaxConnectionPoolSize defines pool size
	MaxConnectionPoolSize int = 50
	// DefaultRequestTimeoutSeconds timeout value
	DefaultRequestTimeoutSeconds int = 60
	// ApiKeyHeaderName header name for API key
	ApiKeyHeaderName string = "X-API-Key"
	// TheQuickBrownFox example
	TheQuickBrownFox string = "fox"
	// ThisIsAVeryLongConstantName long name
	ThisIsAVeryLongConstantName string = "long"
	// DBMaxConnections database connections
	DBMaxConnections int = 100
)

// =============================================================================
// SECTION 8: Status codes and common constants
// =============================================================================

const (
	// StatusOK success status
	StatusOK int = 200
	// StatusCreated resource created
	StatusCreated int = 201
	// StatusAccepted request accepted
	StatusAccepted int = 202
	// StatusNotFound not found
	StatusNotFound int = 404
	// StatusInternalError internal error
	StatusInternalError int = 500
)

// =============================================================================
// SECTION 9: All uppercase without underscores (valid but unusual)
// =============================================================================

const (
	// MAXSIZE all caps no underscore (valid Go)
	MAXSIZE int = 100
	// TIMEOUT all caps no underscore
	TIMEOUT int = 30
	// DEBUG all caps no underscore
	DEBUG bool = true
	// OK all caps no underscore
	OK int = 1
)

// =============================================================================
// SECTION 10: Blank identifier (special case - should be skipped)
// =============================================================================

const (
	// Blank identifier used to ignore value
	_ int = 999
)

// =============================================================================
// SECTION 11: Very long names
// =============================================================================

const (
	// VeryLongConstantNameWithManyWordsInCamelCaseFormat is a long name
	VeryLongConstantNameWithManyWordsInCamelCaseFormat int = 1
	// ThisIsAnExtremelyLongConstantNameThatShouldStillBeValidCamelCase long
	ThisIsAnExtremelyLongConstantNameThatShouldStillBeValidCamelCase int = 2
	// shortName is a short name for comparison
	shortName int = 3
)

// =============================================================================
// SECTION 12: Iota patterns (tests that naming is checked, not iota)
// =============================================================================

const (
	// StatusPending pending status with iota
	StatusPending int = iota
	// StatusActive active status
	StatusActive
	// StatusDone done status
	StatusDone
)

const (
	// FlagA first flag
	FlagA int = 1 << iota
	// FlagB second flag
	FlagB
	// FlagC third flag
	FlagC
)

// =============================================================================
// SECTION 13: Variable declaration (should be ignored by CONST-003)
// =============================================================================

var (
	// variableWithUnderscore should be ignored by const003 rule
	variableWithUnderscore int = MaxSize
	// VARIABLE_CAPS should be ignored by const003 rule
	VARIABLE_CAPS int = MinRetryCount
)

// =============================================================================
// SECTION 14: Type constants (iota with custom type)
// =============================================================================

// Status is a custom status type
type Status int

const (
	// Unknown is unknown status
	Unknown Status = iota
	// Pending is pending status
	Pending
	// Running is running status
	Running
	// Completed is completed status
	Completed
)

// Priority is a custom priority type
type Priority int

const (
	// Low is low priority
	Low Priority = iota
	// Medium is medium priority
	Medium
	// High is high priority
	High
	// Critical is critical priority
	Critical
)
