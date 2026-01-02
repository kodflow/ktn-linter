// Package const004 contains test cases for KTN-CONST-004.
// This file contains ALL cases that MUST trigger KTN-CONST-004 errors.
package const004

// =============================================================================
// SECTION 1: Single character uppercase constants
// =============================================================================

const (
	// A is a single character constant (too short)
	A int = 1 // want "KTN-CONST-004"

	// B is a single character constant (too short)
	B int = 2 // want "KTN-CONST-004"

	// C is a single character constant (too short)
	C string = "c" // want "KTN-CONST-004"

	// D is a single character constant (too short)
	D bool = true // want "KTN-CONST-004"

	// E is a single character constant (too short)
	E float64 = 2.71828 // want "KTN-CONST-004"
)

// =============================================================================
// SECTION 2: Single character lowercase constants
// =============================================================================

const (
	// x is a single character unexported constant (too short)
	x int = 10 // want "KTN-CONST-004"

	// y is a single character unexported constant (too short)
	y int = 20 // want "KTN-CONST-004"

	// z is a single character unexported constant (too short)
	z string = "z" // want "KTN-CONST-004"
)

// =============================================================================
// SECTION 3: Single character in multi-declaration
// =============================================================================

const (
	// M and N are both single character constants (too short)
	M, N int = 1, 2 // want "KTN-CONST-004" "KTN-CONST-004"
)
