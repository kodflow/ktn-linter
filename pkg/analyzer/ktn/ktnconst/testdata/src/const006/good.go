// Package const006 contains test cases for KTN-CONST-006.
// This file contains ALL cases that MUST NOT trigger KTN-CONST-006 errors.
package const006

// =============================================================================
// SECTION 1: Valid constant names (not built-ins)
// =============================================================================

const (
	// MaxSize is a valid constant name
	MaxSize int = 1024

	// MinSize is a valid constant name
	MinSize int = 64

	// DefaultTimeout is a valid constant name
	DefaultTimeout int = 30

	// BufferSize is a valid constant name
	BufferSize int = 4096
)

// =============================================================================
// SECTION 2: Names similar to but not shadowing built-ins
// =============================================================================

const (
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

const (
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

const (
	// INT is valid (Go is case-sensitive, int != INT)
	INT int = 1

	// STRING is valid (string != STRING)
	STRING string = "valid"

	// BOOL is valid (bool != BOOL)
	BOOL bool = true

	// LEN is valid (len != LEN)
	LEN int = 5

	// NIL is valid (nil != NIL) - note: can't assign nil to const
	NIL int = 0
)

// =============================================================================
// SECTION 5: Blank identifier (always allowed)
// =============================================================================

const (
	// Blank identifier is allowed
	_ int = 0

	// Another blank
	_ string = "ignored"
)

// =============================================================================
// SECTION 6: Common valid constant patterns
// =============================================================================

const (
	// StatusOK is a valid status constant
	StatusOK int = 200

	// StatusError is a valid status constant
	StatusError int = 500

	// DefaultPort is a valid network constant
	DefaultPort int = 8080

	// APIVersion is a valid version constant
	APIVersion string = "v1"
)
