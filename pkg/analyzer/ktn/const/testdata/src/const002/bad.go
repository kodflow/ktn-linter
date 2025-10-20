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
