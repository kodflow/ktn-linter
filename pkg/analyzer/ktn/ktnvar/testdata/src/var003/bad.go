// Bad examples for the var004 test case.
package var003

// Constants to avoid magic numbers
const (
	// TEST_AGE is the test age value
	TEST_AGE int = 25
	// TEST_X is the test x coordinate
	TEST_X int = 10
	// TEST_Y is the test y coordinate
	TEST_Y int = 20
	// TEST_ANSWER is the answer value
	TEST_ANSWER int = 42
	// TEST_ONE is the value one
	TEST_ONE int = 1
	// TEST_TWO is the value two
	TEST_TWO int = 2
	// TEST_THREE is the value three
	TEST_THREE int = 3
	// THRESHOLD_VALUE is the threshold for comparison
	THRESHOLD_VALUE int = 5
	// MULTIPLIER is the multiplication factor
	MULTIPLIER int = 2
	// ZERO_INDEX is the starting index
	ZERO_INDEX int = 0
)

// Package-level variables with explicit types (OK, not checked by VAR-005)
var (
	// PackageLevel is a package-level variable
	PackageLevel string = "ok"
	// AnotherGlobal is another package-level variable
	AnotherGlobal int = TEST_ANSWER
)

// badSimpleVar shows incorrect use of var with initialization.
// Should use := instead of var for local variables.
func badSimpleVar() {
	// Variable declarations with initialization (should use :=)
	var name = "Bob"
	// age variable with initialization
	var age = TEST_AGE

	// Using variables to avoid unused warnings
	_ = name
	_ = age
}

// badMultipleVars shows incorrect use of var with multiple variables.
// Should use := instead of var for local variables.
func badMultipleVars() {
	// x variable with initialization
	var x = TEST_X
	// y variable with initialization
	var y = TEST_Y

	// Using variables to avoid unused warnings
	_ = x
	_ = y
}

// badBoolVar shows incorrect use of var with bool.
// Should use := instead of var for local variables.
func badBoolVar() {
	// Boolean variable declaration (should use :=)
	var isEnabled = true

	// Using variable to avoid unused warning
	_ = isEnabled
}

// badStringVar shows incorrect use of var with string.
// Should use := instead of var for local variables.
func badStringVar() {
	// String variable declaration (should use :=)
	var message = "Hello"

	// Using variable to avoid unused warning
	_ = message
}

// badSliceVar shows incorrect use of var with slice.
// Should use := instead of var for local variables.
func badSliceVar() {
	// Slice variable declaration (should use :=)
	var numbers = []int{TEST_ONE, TEST_TWO, TEST_THREE}

	// Using variable to avoid unused warning
	_ = numbers
}

// badMapVar shows incorrect use of var with map.
// Should use := instead of var for local variables.
func badMapVar() {
	// Map variable declaration (should use :=)
	var config = map[string]int{"key": TEST_ONE}

	// Using variable to avoid unused warning
	_ = config
}

// badVarInIf shows incorrect use of var in if block.
// Local variables should use := even inside control structures.
func badVarInIf() {
	// Initialize test variable
	x := TEST_X

	// Check if x exceeds threshold
	if x > THRESHOLD_VALUE {
		// Variable in if block (should use :=)
		var result = "big"

		// Using variable to avoid unused warning
		_ = result
	}
}

// badVarInFor shows incorrect use of var in for loop.
// Local variables should use := even inside loops.
func badVarInFor() {
	// Loop through range
	for i := ZERO_INDEX; i < TEST_THREE; i++ {
		// Variable in for loop (should use :=)
		var item = i * MULTIPLIER

		// Using variable to avoid unused warning
		_ = item
	}
}

// badVarInRange shows incorrect use of var in range loop.
// Local variables should use := even inside range loops.
func badVarInRange() {
	// Initialize slice for iteration
	nums := []int{TEST_ONE, TEST_TWO, TEST_THREE}

	// Iterate over slice
	for _, n := range nums {
		// Variable in range loop (should use :=)
		var doubled = n * MULTIPLIER

		// Using variable to avoid unused warning
		_ = doubled
	}
}

// badVarInSwitch shows incorrect use of var in switch.
// Local variables should use := even in switch cases.
func badVarInSwitch() {
	// Initialize test value
	x := THRESHOLD_VALUE

	// Switch on value
	switch x {
	// Case when x equals threshold
	case THRESHOLD_VALUE:
		// Variable in case block (should use :=)
		var msg = "five"

		// Using variable to avoid unused warning
		_ = msg
	// Default case for other values
	default:
		// Variable in default block (should use :=)
		var other = "other"

		// Using variable to avoid unused warning
		_ = other
	}
}

// badVarInSelect shows incorrect use of var in select.
// Local variables should use := even in select cases.
func badVarInSelect() {
	// Create channels for testing
	ch := make(chan int)

	// Select on channel operations
	select {
	// Case when receiving from channel
	case val := <-ch:
		// Variable in select case (should use :=)
		var result = val * MULTIPLIER

		// Using variable to avoid unused warning
		_ = result
	// Default case when no channel operation is ready
	default:
		// Variable in default block (should use :=)
		var msg = "no data"

		// Using variable to avoid unused warning
		_ = msg
	}
}

// init utilise les fonctions privées
func init() {
	// Appel de badSimpleVar
	badSimpleVar()
	// Appel de badMultipleVars
	badMultipleVars()
	// Appel de badBoolVar
	badBoolVar()
	// Appel de badStringVar
	badStringVar()
	// Appel de badSliceVar
	badSliceVar()
	// Appel de badMapVar
	badMapVar()
	// Appel de badVarInIf
	badVarInIf()
	// Appel de badVarInFor
	badVarInFor()
	// Appel de badVarInRange
	badVarInRange()
	// Appel de badVarInSwitch
	badVarInSwitch()
	// Appel de badVarInSelect
	badVarInSelect()
}
