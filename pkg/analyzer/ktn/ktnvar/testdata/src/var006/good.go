// Package var006 contains test cases for KTN-VAR-006.
// This file contains ALL cases that MUST NOT trigger KTN-VAR-006 errors.
package var006

// =============================================================================
// SECTION 1: Valid variable names (not built-ins)
// =============================================================================

var (
	// MaxSize is a valid variable name
	MaxSize int = 1024

	// MinSize is a valid variable name
	MinSize int = 64

	// DefaultTimeout is a valid variable name
	DefaultTimeout int = 30

	// BufferSize is a valid variable name
	BufferSize int = 4096
)

// =============================================================================
// SECTION 2: Names similar to but not shadowing built-ins
// =============================================================================

var (
	// IsTrue is similar to true but not shadowing
	IsTrue bool = true

	// IsFalse is similar to false but not shadowing
	IsFalse bool = false

	// MyInt is similar to int but not shadowing
	MyInt int = 42

	// MyString is similar to string but not shadowing
	MyString string = "test"

	// MaxLen is similar to len but not shadowing
	MaxLen int = 100

	// NewValue is similar to new but not shadowing
	NewValue int = 0

	// AppendMode is similar to append but not shadowing
	AppendMode int = 1

	// MakeConfig is similar to make but not shadowing
	MakeConfig int = 2
)

// =============================================================================
// SECTION 3: CamelCase versions of built-ins (valid)
// =============================================================================

var (
	// BoolValue is valid (not bool)
	BoolValue bool = true

	// IntValue is valid (not int)
	IntValue int = 42

	// StringValue is valid (not string)
	StringValue string = "valid"

	// ErrorCode is valid (not error)
	ErrorCode int = 1

	// LenValue is valid (not len)
	LenValue int = 10

	// CapValue is valid (not cap)
	CapValue int = 20
)

// =============================================================================
// SECTION 4: Uppercase versions of built-ins (valid in Go)
// =============================================================================

var (
	// INT is valid (Go is case-sensitive, int != INT)
	INT int = 1

	// STRING is valid (string != STRING)
	STRING string = "valid"

	// BOOL is valid (bool != BOOL)
	BOOL bool = true

	// LEN is valid (len != LEN)
	LEN int = 5

	// NIL is valid (nil != NIL)
	NIL int = 0
)

// =============================================================================
// SECTION 5: Blank identifier (always allowed)
// =============================================================================

var (
	// Blank identifier is allowed
	_ int = 0

	// Another blank
	_ string = "ignored"
)

// =============================================================================
// SECTION 6: Common valid variable patterns
// =============================================================================

var (
	// StatusOK is a valid status variable
	StatusOK int = 200

	// StatusError is a valid status variable
	StatusError int = 500

	// DefaultPort is a valid network variable
	DefaultPort int = 8080

	// APIVersion is a valid version variable
	APIVersion string = "v1"
)
