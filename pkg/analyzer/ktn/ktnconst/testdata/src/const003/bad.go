// Package const003 contains test cases for KTN-CONST-003.
// This file contains ALL cases that MUST trigger KTN-CONST-003 errors.
// Go convention: Use CamelCase (MixedCaps), not underscores.
package const003

// =============================================================================
// SECTION 1: SCREAMING_SNAKE_CASE (C/Java convention, NOT Go)
// =============================================================================

const (
	// MaxSizeScream contains underscore
	MAX_SIZE int = 100 // want "KTN-CONST-003"
	// ApiKeyScream contains underscore
	API_KEY string = "secret" // want "KTN-CONST-003"
	// HttpTimeoutScream contains underscore
	HTTP_TIMEOUT int = 30 // want "KTN-CONST-003"
	// DatabaseMaxConnections contains underscore
	DB_MAX_CONNECTIONS int = 100 // want "KTN-CONST-003"
	// DefaultPortNumber contains underscore
	DEFAULT_PORT int = 8080 // want "KTN-CONST-003"
	// IsProductionMode contains underscore
	IS_PRODUCTION bool = false // want "KTN-CONST-003"
)

// =============================================================================
// SECTION 2: snake_case (lowercase with underscores)
// =============================================================================

const (
	// maxSizeSnake contains underscore
	max_size int = 100 // want "KTN-CONST-003"
	// apiKeySnake contains underscore
	api_key string = "secret" // want "KTN-CONST-003"
	// httpTimeoutSnake contains underscore
	http_timeout int = 30 // want "KTN-CONST-003"
	// defaultPortSnake contains underscore
	default_port int = 8080 // want "KTN-CONST-003"
	// isEnabledSnake contains underscore
	is_enabled bool = true // want "KTN-CONST-003"
)

// =============================================================================
// SECTION 3: Mixed_Case (PascalCase with underscores)
// =============================================================================

const (
	// MaxSizeMixed contains underscore
	Max_Size int = 100 // want "KTN-CONST-003"
	// ApiKeyMixed contains underscore
	Api_Key string = "secret" // want "KTN-CONST-003"
	// HttpTimeoutMixed contains underscore
	Http_Timeout int = 30 // want "KTN-CONST-003"
	// DefaultConfigMixed contains underscore
	Default_Config string = "cfg" // want "KTN-CONST-003"
)

// =============================================================================
// SECTION 4: Minimal underscore cases (single underscore in middle)
// =============================================================================

const (
	// TwoLettersUnderscore single underscore between letters
	a_b int = 1 // want "KTN-CONST-003"
	// ThreeLettersUnderscore single underscore in middle
	ab_c int = 2 // want "KTN-CONST-003"
	// LongerUnderscore single underscore
	abc_def int = 3 // want "KTN-CONST-003"
	// UpperSingleUnderscore single underscore
	A_B int = 4 // want "KTN-CONST-003"
	// MixedSingleUnderscore single underscore
	Ab_Cd int = 5 // want "KTN-CONST-003"
)

// =============================================================================
// SECTION 5: Multiple underscores
// =============================================================================

const (
	// DoubleUnderscore contains two consecutive underscores
	double__underscore int = 1 // want "KTN-CONST-003"
	// TripleUnderscore contains three consecutive underscores
	triple___underscore int = 2 // want "KTN-CONST-003"
	// ManyUnderscores contains many underscores
	a___b___c int = 3 // want "KTN-CONST-003"
	// ComplexMultiple complex pattern
	MAX_CONNECTION_POOL_SIZE int = 50 // want "KTN-CONST-003"
	// VeryLongScream very long screaming name
	DEFAULT_REQUEST_TIMEOUT_SECONDS int = 60 // want "KTN-CONST-003"
)

// =============================================================================
// SECTION 6: Edge cases - Position of underscore
// =============================================================================

const (
	// LeadingUnderscore starts with underscore (not blank identifier)
	_leading int = 1 // want "KTN-CONST-003"
	// TrailingUnderscore ends with underscore
	trailing_ int = 2 // want "KTN-CONST-003"
	// BothEndsUnderscore underscore at both ends
	_both_ends_ int = 3 // want "KTN-CONST-003"
	// LeadingUnderscoreUpper starts with underscore then upper
	_Leading int = 4 // want "KTN-CONST-003"
	// TrailingUnderscoreNum ends with underscore before value
	value_ int = 5 // want "KTN-CONST-003"
)

// =============================================================================
// SECTION 7: Underscores with numbers
// =============================================================================

const (
	// UnderscoreBetweenNumbers underscore between digits
	value_1_2 int = 1 // want "KTN-CONST-003"
	// UnderscoreAfterNumber underscore after number
	http2_server string = "h2" // want "KTN-CONST-003"
	// UnderscoreBeforeNumber underscore before number
	version_2 int = 2 // want "KTN-CONST-003"
	// NumberUnderscoreNumber number underscore number
	v1_0_0 string = "1.0.0" // want "KTN-CONST-003"
	// ComplexNumberPattern complex with numbers
	api_v2_endpoint string = "/v2" // want "KTN-CONST-003"
)

// =============================================================================
// SECTION 8: Long multi-word snake_case
// =============================================================================

const (
	// ThreeWordSnake three words with underscores
	the_quick_fox string = "fox" // want "KTN-CONST-003"
	// FourWordSnake four words with underscores
	the_quick_brown_fox string = "brown" // want "KTN-CONST-003"
	// LongPhraseSnake long phrase
	this_is_a_very_long_constant_name string = "long" // want "KTN-CONST-003"
	// AllCapsLong all caps long
	API_KEY_HEADER_NAME_VALUE string = "X-API" // want "KTN-CONST-003"
)

// =============================================================================
// SECTION 9: Acronyms with underscores (incorrect style)
// =============================================================================

const (
	// HttpApiUnderscore HTTP with underscore
	HTTP_API string = "/api" // want "KTN-CONST-003"
	// UrlPathUnderscore URL with underscore
	URL_PATH string = "/path" // want "KTN-CONST-003"
	// XmlHttpUnderscore XML HTTP with underscore
	XML_HTTP_REQUEST string = "xhr" // want "KTN-CONST-003"
	// AppIdUnderscore App ID with underscore
	APP_ID string = "id123" // want "KTN-CONST-003"
	// TcpIpUnderscore TCP IP with underscore
	TCP_IP_PROTOCOL string = "tcp" // want "KTN-CONST-003"
)
