// Package const004 contains test cases for KTN-CONST-004.
// This file contains ALL cases that MUST NOT trigger KTN-CONST-004 errors.
package const004

// =============================================================================
// SECTION 1: Valid names with 2+ characters
// =============================================================================

const (
	// GoodMinLen is a constant with minimum valid length (2 chars)
	GoodMinLen int = 42

	// GoodThreeChars has three characters
	GoodThreeChars int = 100

	// GoodLongName has a descriptive name
	GoodLongName string = "hello"

	// GoodCamelCase is a proper CamelCase name
	GoodCamelCase bool = true

	// GoodWithNumbers uses letters and numbers
	GoodWithNumbers int = 123

	// GoodMaxLength is a reasonably long name
	GoodMaxLength float64 = 3.14159
)

// =============================================================================
// SECTION 2: Two character names (minimum valid)
// =============================================================================

const (
	// ID is a common two-letter constant
	ID int = 1

	// OK is a two-letter constant
	OK int = 0

	// IP is a network-related constant
	IP string = "127.0.0.1"

	// US is a country code constant
	US string = "US"

	// GB is a country code constant
	GB string = "GB"
)

// =============================================================================
// SECTION 3: Blank identifier (always allowed)
// =============================================================================

const (
	// Blank identifier is intentionally ignored
	_ int = 0

	// Another blank identifier
	_ string = "ignored"

	// Blank with iota
	_ int = iota
)

// =============================================================================
// SECTION 4: Exported constants with proper names
// =============================================================================

const (
	// MaxSize is an exported constant
	MaxSize int = 1024

	// MinSize is an exported constant
	MinSize int = 64

	// DefaultTimeout is an exported constant
	DefaultTimeout int = 30

	// APIVersion is an exported constant
	APIVersion string = "v1"
)

// =============================================================================
// SECTION 5: Unexported constants with proper names
// =============================================================================

const (
	// maxRetries is an unexported constant
	maxRetries int = 3

	// defaultPort is an unexported constant
	defaultPort int = 8080

	// bufferSize is an unexported constant
	bufferSize int = 4096
)

// =============================================================================
// SECTION 6: Iota with proper names
// =============================================================================

const (
	// StatusOK is the first status
	StatusOK int = iota

	// StatusError is an error status
	StatusError

	// StatusPending is a pending status
	StatusPending
)

// =============================================================================
// SECTION 7: Multi-name declarations with valid names
// =============================================================================

const (
	// MultiA and MultiB are valid multi-declaration names
	MultiA, MultiB int = 1, 2

	// StrA and StrB are valid string constants
	StrA, StrB string = "a", "b"
)
