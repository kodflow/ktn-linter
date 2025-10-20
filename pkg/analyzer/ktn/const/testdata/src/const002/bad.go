package const002

// Bad: Multiple const blocks (scattered) - violates KTN-CONST-002

// First group of constants
const (
	// FIRST_GROUP_A is in the first group
	FIRST_GROUP_A string = "first"
	// FIRST_GROUP_B is in the first group
	FIRST_GROUP_B string = "group"
)

// Second group - scattered (ERROR #1)
const (
	// SECOND_GROUP is in a second group
	SECOND_GROUP string = "scattered"
)

// Variable declaration
var GlobalVar string = "test"

// BAD_AFTER_VAR appears after var (ERROR #2)
const (
	// BAD_AFTER_VAR is placed after vars
	BAD_AFTER_VAR string = "bad"
)

// Another scattered const after var (ERROR #3)
const (
	// SCATTERED_ONE in another block
	SCATTERED_ONE string = "one"
)

// Yet another scattered const after var (ERROR #4)
const (
	// SCATTERED_TWO in yet another block
	SCATTERED_TWO string = "two"
)

// Another const after var (ERROR #5)
const (
	// THIRD_AFTER_VAR also after var
	THIRD_AFTER_VAR int = 1
)

// helperFunction is used to demonstrate const blocks separated by other declarations.
//
// Returns:
//   - string: a helper message
func helperFunction() string {
	// Retour de la fonction
	return "helper"
}

// Last const after var (ERROR #6)
const (
	// FOURTH_AFTER_VAR also after var
	FOURTH_AFTER_VAR int = 2
)
