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

// Third scattered const group (ERROR #2)
const (
	// BAD_SCATTERED_A is scattered
	BAD_SCATTERED_A string = "bad"
)

// Fourth scattered const group (ERROR #3)
const (
	// SCATTERED_ONE in another block
	SCATTERED_ONE string = "one"
)

// Fifth scattered const group (ERROR #4)
const (
	// SCATTERED_TWO in yet another block
	SCATTERED_TWO string = "two"
)

// Sixth scattered const group (ERROR #5)
const (
	// THIRD_SCATTERED also scattered
	THIRD_SCATTERED int = 1
)

// Seventh scattered const group (ERROR #6)
const (
	// FOURTH_SCATTERED also scattered
	FOURTH_SCATTERED int = 2
)

// Eighth scattered const group (ERROR #7)
const (
	// CONST_FINAL is yet another scattered const
	CONST_FINAL string = "final"
)

// helperFunction is used to demonstrate const blocks separated by other declarations.
//
// Returns:
//   - string: a helper message
func helperFunction() string {
	// Retour de la fonction
	return "helper"
}
