// Package const005 contains test cases for KTN-CONST-005.
// This file contains ALL cases that MUST trigger KTN-CONST-005 errors.
package const005

// =============================================================================
// SECTION 1: Exported constants exceeding 30 characters
// =============================================================================

const (
	// ThisConstantNameExceedsThirtyCharacters has 40 chars (too long)
	ThisConstantNameExceedsThirtyCharacters int = 1 // want "KTN-CONST-005"

	// VeryLongConstantNameThatExceedsMaximumLength has 45 chars (too long)
	VeryLongConstantNameThatExceedsMaximumLength int = 2 // want "KTN-CONST-005"

	// DefaultHTTPConnectionTimeoutInSeconds has 38 chars (too long)
	DefaultHTTPConnectionTimeoutInSeconds int = 30 // want "KTN-CONST-005"
)

// =============================================================================
// SECTION 2: Unexported constants exceeding 30 characters
// =============================================================================

const (
	// thisUnexportedConstantNameIsTooLong has 36 chars (too long)
	thisUnexportedConstantNameIsTooLong int = 100 // want "KTN-CONST-005"

	// maximumNumberOfRetriesBeforeError has 34 chars (too long)
	maximumNumberOfRetriesBeforeError int = 5 // want "KTN-CONST-005"
)
