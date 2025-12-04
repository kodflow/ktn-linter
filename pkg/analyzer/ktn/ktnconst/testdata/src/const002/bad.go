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
const (
	// BAD_SECOND_GROUP is in a second group
	BAD_SECOND_GROUP string = "scattered" // want "KTN-CONST-002: les constantes doivent être groupées ensemble dans un seul bloc"
)

// Third scattered const group (ERROR #2)
const (
	// BAD_SCATTERED_A is scattered
	BAD_SCATTERED_A string = "bad" // want "KTN-CONST-002: les constantes doivent être groupées ensemble dans un seul bloc"
)

// Fourth scattered const group (ERROR #3)
const (
	// BAD_SCATTERED_ONE in another block
	BAD_SCATTERED_ONE string = "one" // want "KTN-CONST-002: les constantes doivent être groupées ensemble dans un seul bloc"
)

// Fifth scattered const group (ERROR #4)
const (
	// BAD_SCATTERED_TWO in yet another block
	BAD_SCATTERED_TWO string = "two" // want "KTN-CONST-002: les constantes doivent être groupées ensemble dans un seul bloc"
)

// Sixth scattered const group (ERROR #5)
const (
	// BAD_THIRD_SCATTERED also scattered
	BAD_THIRD_SCATTERED int = 1 // want "KTN-CONST-002: les constantes doivent être groupées ensemble dans un seul bloc"
)

// Seventh scattered const group (ERROR #6)
const (
	// BAD_FOURTH_SCATTERED also scattered
	BAD_FOURTH_SCATTERED int = 2 // want "KTN-CONST-002: les constantes doivent être groupées ensemble dans un seul bloc"
)

// Eighth scattered const group (ERROR #7)
const (
	// BAD_CONST_FINAL is yet another scattered const
	BAD_CONST_FINAL string = "final" // want "KTN-CONST-002: les constantes doivent être groupées ensemble dans un seul bloc"
)

// Ninth scattered const group (ERROR #8)
const (
	// BAD_CONST_EXTRA yet another scattered const block
	BAD_CONST_EXTRA string = "extra" // want "KTN-CONST-002: les constantes doivent être groupées ensemble dans un seul bloc"
)

// Variable declaration (after all constants)
var badGlobalVar string = "some var"
