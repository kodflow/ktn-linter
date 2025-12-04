// Bad examples for the const003 test case.
package const003

// Bad: Invalid naming (violates KTN-CONST-003)
// But respects: explicit types, comments, and single-block grouping
const (
	// camelCase naming (INVALID)
	// badMaxSize in camelCase
	badMaxSize int = 100 // want "KTN-CONST-003: la constante 'badMaxSize' doit utiliser la convention CAPITAL_UNDERSCORE"
	// badApiKey in camelCase
	badApiKey string = "secret" // want "KTN-CONST-003: la constante 'badApiKey' doit utiliser la convention CAPITAL_UNDERSCORE"
	// badHttpTimeout in camelCase
	badHttpTimeout int = 30 // want "KTN-CONST-003: la constante 'badHttpTimeout' doit utiliser la convention CAPITAL_UNDERSCORE"

	// PascalCase naming (INVALID)
	// BadMaxSizePascal in PascalCase
	BadMaxSizePascal int = 100 // want "KTN-CONST-003: la constante 'BadMaxSizePascal' doit utiliser la convention CAPITAL_UNDERSCORE"
	// BadApiKeyPascal in PascalCase
	BadApiKeyPascal string = "secret" // want "KTN-CONST-003: la constante 'BadApiKeyPascal' doit utiliser la convention CAPITAL_UNDERSCORE"
	// BadHttpTimeoutPascal in PascalCase
	BadHttpTimeoutPascal int = 30 // want "KTN-CONST-003: la constante 'BadHttpTimeoutPascal' doit utiliser la convention CAPITAL_UNDERSCORE"

	// snake_case (lowercase with underscores) - INVALID
	// bad_max_size in snake_case
	bad_max_size int = 100 // want "KTN-CONST-003: la constante 'bad_max_size' doit utiliser la convention CAPITAL_UNDERSCORE"
	// bad_api_key in snake_case
	bad_api_key string = "secret" // want "KTN-CONST-003: la constante 'bad_api_key' doit utiliser la convention CAPITAL_UNDERSCORE"
	// bad_http_timeout in snake_case
	bad_http_timeout int = 30 // want "KTN-CONST-003: la constante 'bad_http_timeout' doit utiliser la convention CAPITAL_UNDERSCORE"

	// Mixed case with underscores - INVALID
	// Bad_Max_Size mixed case
	Bad_Max_Size int = 100 // want "KTN-CONST-003: la constante 'Bad_Max_Size' doit utiliser la convention CAPITAL_UNDERSCORE"
	// Bad_Api_Key mixed case
	Bad_Api_Key string = "secret" // want "KTN-CONST-003: la constante 'Bad_Api_Key' doit utiliser la convention CAPITAL_UNDERSCORE"
	// Bad_Http_Timeout mixed case
	Bad_Http_Timeout int = 30 // want "KTN-CONST-003: la constante 'Bad_Http_Timeout' doit utiliser la convention CAPITAL_UNDERSCORE"

	// More camelCase examples
	// badStatusOk in camelCase
	badStatusOk int = 200 // want "KTN-CONST-003: la constante 'badStatusOk' doit utiliser la convention CAPITAL_UNDERSCORE"
	// badStatusCreated in camelCase
	badStatusCreated int = 201 // want "KTN-CONST-003: la constante 'badStatusCreated' doit utiliser la convention CAPITAL_UNDERSCORE"
	// badStatusError in camelCase
	badStatusError int = 500 // want "KTN-CONST-003: la constante 'badStatusError' doit utiliser la convention CAPITAL_UNDERSCORE"

	// More PascalCase examples
	// BadStateIdle in PascalCase
	BadStateIdle int = 0 // want "KTN-CONST-003: la constante 'BadStateIdle' doit utiliser la convention CAPITAL_UNDERSCORE"
	// BadStateRunning in PascalCase
	BadStateRunning int = 1 // want "KTN-CONST-003: la constante 'BadStateRunning' doit utiliser la convention CAPITAL_UNDERSCORE"
	// BadStatePaused in PascalCase
	BadStatePaused int = 2 // want "KTN-CONST-003: la constante 'BadStatePaused' doit utiliser la convention CAPITAL_UNDERSCORE"

	// Mixed variations
	// BadErrorNotFound PascalCase
	BadErrorNotFound string = "not found" // want "KTN-CONST-003: la constante 'BadErrorNotFound' doit utiliser la convention CAPITAL_UNDERSCORE"
	// badErrorUnauthorized camelCase
	badErrorUnauthorized string = "unauthorized" // want "KTN-CONST-003: la constante 'badErrorUnauthorized' doit utiliser la convention CAPITAL_UNDERSCORE"
	// Bad_Error_Internal mixed
	Bad_Error_Internal string = "internal" // want "KTN-CONST-003: la constante 'Bad_Error_Internal' doit utiliser la convention CAPITAL_UNDERSCORE"

	// Starting with lowercase
	// badDefaultPort lowercase start
	badDefaultPort int = 8080 // want "KTN-CONST-003: la constante 'badDefaultPort' doit utiliser la convention CAPITAL_UNDERSCORE"
	// badDefaultHost lowercase start
	badDefaultHost string = "localhost" // want "KTN-CONST-003: la constante 'badDefaultHost' doit utiliser la convention CAPITAL_UNDERSCORE"
	// badDefaultProtocol lowercase start
	badDefaultProtocol string = "http" // want "KTN-CONST-003: la constante 'badDefaultProtocol' doit utiliser la convention CAPITAL_UNDERSCORE"

	// Complex camelCase
	// badMaxConnectionPoolSize complex camelCase
	badMaxConnectionPoolSize int = 50 // want "KTN-CONST-003: la constante 'badMaxConnectionPoolSize' doit utiliser la convention CAPITAL_UNDERSCORE"
	// badDefaultRequestTimeout complex camelCase
	badDefaultRequestTimeout int = 60 // want "KTN-CONST-003: la constante 'badDefaultRequestTimeout' doit utiliser la convention CAPITAL_UNDERSCORE"
	// badApiKeyHeaderName complex camelCase
	badApiKeyHeaderName string = "X-API-Key" // want "KTN-CONST-003: la constante 'badApiKeyHeaderName' doit utiliser la convention CAPITAL_UNDERSCORE"

	// Partially correct (mixed) - INVALID
	// BAD_Max_Size partially correct
	BAD_Max_Size int = 100 // want "KTN-CONST-003: la constante 'BAD_Max_Size' doit utiliser la convention CAPITAL_UNDERSCORE"
	// Bad_Api_KEY partially correct
	Bad_Api_KEY string = "key" // want "KTN-CONST-003: la constante 'Bad_Api_KEY' doit utiliser la convention CAPITAL_UNDERSCORE"
	// BAD_HTTP_timeout partially correct
	BAD_HTTP_timeout int = 30 // want "KTN-CONST-003: la constante 'BAD_HTTP_timeout' doit utiliser la convention CAPITAL_UNDERSCORE"

	// Database constants with wrong naming
	// badDbMaxConnections database setting
	badDbMaxConnections int = 100 // want "KTN-CONST-003: la constante 'badDbMaxConnections' doit utiliser la convention CAPITAL_UNDERSCORE"
	// BadDbMinConnections database setting
	BadDbMinConnections int = 10 // want "KTN-CONST-003: la constante 'BadDbMinConnections' doit utiliser la convention CAPITAL_UNDERSCORE"
	// bad_db_timeout database setting
	bad_db_timeout int = 30 // want "KTN-CONST-003: la constante 'bad_db_timeout' doit utiliser la convention CAPITAL_UNDERSCORE"
)
