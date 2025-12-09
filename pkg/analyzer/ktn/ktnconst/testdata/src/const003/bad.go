// Bad examples for the const003 test case.
package const003

// Bad: Invalid naming (violates KTN-CONST-003)
// But respects: explicit types, comments, and single-block grouping
const (
	// camelCase naming (INVALID)
	// badMaxSize in camelCase
	badMaxSize int = 100 // want "KTN-CONST-003"
	// badApiKey in camelCase
	badApiKey string = "secret" // want "KTN-CONST-003"
	// badHttpTimeout in camelCase
	badHttpTimeout int = 30 // want "KTN-CONST-003"

	// PascalCase naming (INVALID)
	// BadMaxSizePascal in PascalCase
	BadMaxSizePascal int = 100 // want "KTN-CONST-003"
	// BadApiKeyPascal in PascalCase
	BadApiKeyPascal string = "secret" // want "KTN-CONST-003"
	// BadHttpTimeoutPascal in PascalCase
	BadHttpTimeoutPascal int = 30 // want "KTN-CONST-003"

	// snake_case (lowercase with underscores) - INVALID
	// bad_max_size in snake_case
	bad_max_size int = 100 // want "KTN-CONST-003"
	// bad_api_key in snake_case
	bad_api_key string = "secret" // want "KTN-CONST-003"
	// bad_http_timeout in snake_case
	bad_http_timeout int = 30 // want "KTN-CONST-003"

	// Mixed case with underscores - INVALID
	// Bad_Max_Size mixed case
	Bad_Max_Size int = 100 // want "KTN-CONST-003"
	// Bad_Api_Key mixed case
	Bad_Api_Key string = "secret" // want "KTN-CONST-003"
	// Bad_Http_Timeout mixed case
	Bad_Http_Timeout int = 30 // want "KTN-CONST-003"

	// More camelCase examples
	// badStatusOk in camelCase
	badStatusOk int = 200 // want "KTN-CONST-003"
	// badStatusCreated in camelCase
	badStatusCreated int = 201 // want "KTN-CONST-003"
	// badStatusError in camelCase
	badStatusError int = 500 // want "KTN-CONST-003"

	// More PascalCase examples
	// BadStateIdle in PascalCase
	BadStateIdle int = 0 // want "KTN-CONST-003"
	// BadStateRunning in PascalCase
	BadStateRunning int = 1 // want "KTN-CONST-003"
	// BadStatePaused in PascalCase
	BadStatePaused int = 2 // want "KTN-CONST-003"

	// Mixed variations
	// BadErrorNotFound PascalCase
	BadErrorNotFound string = "not found" // want "KTN-CONST-003"
	// badErrorUnauthorized camelCase
	badErrorUnauthorized string = "unauth" // want "KTN-CONST-003"
	// Bad_Error_Internal mixed
	Bad_Error_Internal string = "internal" // want "KTN-CONST-003"

	// Starting with lowercase
	// badDefaultPort lowercase start
	badDefaultPort int = 8080 // want "KTN-CONST-003"
	// badDefaultHost lowercase start
	badDefaultHost string = "localhost" // want "KTN-CONST-003"
	// badDefaultProtocol lowercase start
	badDefaultProtocol string = "http" // want "KTN-CONST-003"

	// Complex camelCase
	// badMaxConnectionPoolSize complex camelCase
	badMaxConnectionPoolSize int = 50 // want "KTN-CONST-003"
	// badDefaultRequestTimeout complex camelCase
	badDefaultRequestTimeout int = 60 // want "KTN-CONST-003"
	// badApiKeyHeaderName complex camelCase
	badApiKeyHeaderName string = "X-Key" // want "KTN-CONST-003"

	// Partially correct (mixed) - INVALID
	// BAD_Max_Size partially correct
	BAD_Max_Size int = 100 // want "KTN-CONST-003"
	// Bad_Api_KEY partially correct
	Bad_Api_KEY string = "key" // want "KTN-CONST-003"
	// BAD_HTTP_timeout partially correct
	BAD_HTTP_timeout int = 30 // want "KTN-CONST-003"

	// Database constants with wrong naming
	// badDbMaxConnections database setting
	badDbMaxConnections int = 100 // want "KTN-CONST-003"
	// BadDbMinConnections database setting
	BadDbMinConnections int = 10 // want "KTN-CONST-003"
	// bad_db_timeout database setting
	bad_db_timeout int = 30 // want "KTN-CONST-003"
)
