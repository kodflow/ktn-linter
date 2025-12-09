// Bad examples for the const002 test case.
package const002

// Bad: Multiple scattered const blocks (violates KTN-CONST-002)

// First group of constants (OK - first block)
const (
	// BAD_FIRST_GROUP_A is in the first group
	BAD_FIRST_GROUP_A string = "first"
	// BAD_FIRST_GROUP_B is in the first group
	BAD_FIRST_GROUP_B string = "group"
)

// Second scattered const group (ERROR #1)
const ( // want "KTN-CONST-002"
	// BAD_SECOND_GROUP is in a second group
	BAD_SECOND_GROUP string = "scattered"
)

// Third scattered const group (ERROR #2)
const ( // want "KTN-CONST-002"
	// BAD_SCATTERED_A is scattered
	BAD_SCATTERED_A string = "bad"
)

// Fourth scattered const group (ERROR #3)
const ( // want "KTN-CONST-002"
	// BAD_SCATTERED_ONE in another block
	BAD_SCATTERED_ONE string = "one"
)

// Fifth scattered const group (ERROR #4)
const ( // want "KTN-CONST-002"
	// BAD_SCATTERED_TWO in yet another block
	BAD_SCATTERED_TWO string = "two"
)

// Sixth scattered const group (ERROR #5)
const ( // want "KTN-CONST-002"
	// BAD_THIRD_SCATTERED also scattered
	BAD_THIRD_SCATTERED int = 1
)

// Seventh scattered const group (ERROR #6)
const ( // want "KTN-CONST-002"
	// BAD_FOURTH_SCATTERED also scattered
	BAD_FOURTH_SCATTERED int = 2
)

// Eighth scattered const group (ERROR #7)
const ( // want "KTN-CONST-002"
	// BAD_CONST_FINAL is yet another scattered const
	BAD_CONST_FINAL string = "final"
)

// Ninth scattered const group (ERROR #8)
const ( // want "KTN-CONST-002"
	// BAD_CONST_EXTRA yet another scattered const block
	BAD_CONST_EXTRA string = "extra"
)

// Variable declaration (after all constants)
var badGlobalVar string = "some var"
