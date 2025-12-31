// Package const003 provides good test cases.
package const003

// Good: All constants use CamelCase naming (Go standard)
const (
	// === Single letter constants (valid) ===

	// A represents the first value
	A int = 1
	// B represents the second value
	B int = 2
	// C represents the third value
	C int = 3

	// === PascalCase for exported constants ===

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

	// === Acronyms in PascalCase ===

	// APIEndpoint is the API endpoint
	APIEndpoint string = "/api"
	// HTTPStatus is the HTTP status code
	HTTPStatus int = 200
	// URLPath is the URL path
	URLPath string = "/path"

	// === camelCase for unexported constants ===

	// maxInternalSize is internal max size
	maxInternalSize int = 50
	// defaultTimeout is the default timeout
	defaultTimeout int = 10
	// httpInternalPort is internal port
	httpInternalPort int = 8080

	// === Constants with numbers ===

	// Http2Protocol is HTTP/2 protocol
	Http2Protocol string = "h2"
	// Tls12Version is TLS 1.2 version
	Tls12Version string = "1.2"
	// Version100 is version 1.0.0
	Version100 string = "1.0.0"

	// === Complex names without underscores ===

	// MaxConnectionPoolSize defines pool size
	MaxConnectionPoolSize int = 50
	// DefaultRequestTimeoutSeconds timeout value
	DefaultRequestTimeoutSeconds int = 60
	// ApiKeyHeaderName header name for API key
	ApiKeyHeaderName string = "X-API-Key"

	// === Status codes ===

	// StatusOK success status
	StatusOK int = 200
	// StatusCreated resource created
	StatusCreated int = 201
	// StatusAccepted request accepted
	StatusAccepted int = 202

	// === Special cases ===

	// Blank identifier (valid - should be skipped)
	_ int = 999
	// TestValue used for testing
	TestValue int = 123

	// === Constant with explicit type (T3.2) ===
	// Testing CONST-003 CamelCase naming with explicit types

	// TypedIntValue is a constant with explicit int type
	TypedIntValue int = 42
	// TypedStringValue is a string with explicit type
	TypedStringValue string = "test"

	// === Very long names (T3.3) ===

	// VeryLongConstantNameWithManyWordsInCamelCaseFormat is a long name
	VeryLongConstantNameWithManyWordsInCamelCaseFormat int = 1
	// ThisIsAnExtremelyLongConstantNameThatShouldStillBeValidCamelCase long
	ThisIsAnExtremelyLongConstantNameThatShouldStillBeValidCamelCase int = 2
	// shortName is a short name for comparison
	shortName int = 3
)

// Variable declaration to test that only const are checked
var (
	// variableNotConst should be ignored by const003 rule
	variableNotConst int = TestValue
)
