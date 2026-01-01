// Package const005 contains test cases for KTN-CONST-005.
// This file contains ALL cases that MUST NOT trigger KTN-CONST-005 errors.
package const005

// =============================================================================
// SECTION 1: Valid names with <= 30 characters
// =============================================================================

const (
	// ShortName is a short constant name
	ShortName int = 1

	// MediumLengthConstantName has 24 chars
	MediumLengthConstantName int = 2

	// MaxLengthConstantNameHere has exactly 30 characters (maximum allowed)
	MaxLengthConstantNameHere int = 3
)

// =============================================================================
// SECTION 2: Edge cases at boundary (30 chars)
// =============================================================================

const (
	// ThisConstantNameHas30Chars has exactly 30 characters
	ThisConstantNameHas30Chars int = 30

	// ExactlyThirtyCharactersNa is 30 chars
	ExactlyThirtyCharactersNa string = "30"
)

// =============================================================================
// SECTION 3: Common naming patterns under limit
// =============================================================================

const (
	// DefaultHTTPTimeoutSeconds is a typical constant name
	DefaultHTTPTimeoutSeconds int = 30

	// MaxConcurrentConnections is a typical constant
	MaxConcurrentConnections int = 100

	// MinBufferSizeBytes is a typical constant
	MinBufferSizeBytes int = 1024

	// APIVersionNumber is a version constant
	APIVersionNumber string = "v2"
)

// =============================================================================
// SECTION 4: Short constants
// =============================================================================

const (
	// ID is a two-letter constant
	ID int = 1

	// OK is a short constant
	OK int = 0

	// Pi is a mathematical constant
	Pi float64 = 3.14159
)

// =============================================================================
// SECTION 5: Blank identifier (always allowed)
// =============================================================================

const (
	// Blank identifier is allowed
	_ int = 0
)

// =============================================================================
// SECTION 6: Unexported constants under limit
// =============================================================================

const (
	// maxRetriesBeforeFailover is 24 chars
	maxRetriesBeforeFailover int = 3

	// defaultConnectionTimeout is 24 chars
	defaultConnectionTimeout int = 5
)
