package const002

// Bad: Multiple const blocks (scattered) - violates KTN-CONST-002
// But respects other rules: explicit types, proper naming, and comments

// First group of constants
const (
	// FIRST_GROUP_A is in the first group
	FIRST_GROUP_A string = "first"
	// FIRST_GROUP_B is in the first group
	FIRST_GROUP_B string = "group"
)

// SECOND_GROUP is scattered
const (
	// SECOND_GROUP is in a second group
	SECOND_GROUP string = "scattered"
)

// Variable declaration
var GlobalVar string = "test"

// BAD_AFTER_VAR appears after var
const (
	// BAD_AFTER_VAR is placed after vars
	BAD_AFTER_VAR string = "bad"
)

// Edge case: Another scattered const (no var between these)
const (
	// SCATTERED_ONE in another block
	SCATTERED_ONE string = "one"
)

const (
	// SCATTERED_TWO in yet another block
	SCATTERED_TWO string = "two"
)

// Edge case: Scattered consts WITHOUT any var declarations
const (
	// FIRST_NO_VAR in first block
	FIRST_NO_VAR int = 1
)

// Another function to separate const blocks (tests non-GenDecl)
func separatorFunc() {}

const (
	// SECOND_NO_VAR in second block (scattered, no var)
	SECOND_NO_VAR int = 2
)
