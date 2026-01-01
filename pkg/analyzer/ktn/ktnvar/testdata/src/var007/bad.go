// Package var003 contains test cases for KTN rules.
package var003

// Constants to avoid magic numbers
const (
	// TestAge is the test age value
	TestAge int = 25
	// TestX is the test x coordinate
	TestX int = 10
	// TestY is the test y coordinate
	TestY int = 20
	// TestAnswer is the answer value
	TestAnswer int = 42
	// TestOne is the value one
	TestOne int = 1
	// TestTwo is the value two
	TestTwo int = 2
	// TestThree is the value three
	TestThree int = 3
	// ThresholdValue is the threshold for comparison
	ThresholdValue int = 5
	// MultiplierValue is the multiplication factor
	MultiplierValue int = 2
	// ZeroIndex is the starting index
	ZeroIndex int = 0
)

// Package-level variables with explicit types (OK, not checked by VAR-007)
var (
	// PackageLevel is a package-level variable
	PackageLevel string = "ok"
	// AnotherGlobal is another package-level variable
	AnotherGlobal int = TestAnswer
)

// badSimpleVar shows incorrect use of var with initialization.
// Should use := instead of var for local variables.
func badSimpleVar() {
	// Variable declarations with initialization (should use :=)
	var name = "Bob" // want "KTN-VAR-007"
	// age variable with initialization
	var age = TestAge // want "KTN-VAR-007"

	// Using variables to avoid unused warnings
	_ = name
	_ = age
}

// badMultipleVars shows incorrect use of var with multiple variables.
// Should use := instead of var for local variables.
func badMultipleVars() {
	// x variable with initialization
	var x = TestX // want "KTN-VAR-007"
	// y variable with initialization
	var y = TestY // want "KTN-VAR-007"

	// Using variables to avoid unused warnings
	_ = x
	_ = y
}

// badBoolVar shows incorrect use of var with bool.
// Should use := instead of var for local variables.
func badBoolVar() {
	// Boolean variable declaration (should use :=)
	var isEnabled = true // want "KTN-VAR-007"

	// Using variable to avoid unused warning
	_ = isEnabled
}

// badStringVar shows incorrect use of var with string.
// Should use := instead of var for local variables.
func badStringVar() {
	// String variable declaration (should use :=)
	var message = "Hello" // want "KTN-VAR-007"

	// Using variable to avoid unused warning
	_ = message
}

// badSliceVar shows incorrect use of var with slice.
// Should use := instead of var for local variables.
func badSliceVar() {
	// Slice variable declaration (should use :=)
	var numbers = []int{TestOne, TestTwo, TestThree} // want "KTN-VAR-007"

	// Using variable to avoid unused warning
	_ = numbers
}

// badMapVar shows incorrect use of var with map.
// Should use := instead of var for local variables.
func badMapVar() {
	// Map variable declaration (should use :=)
	var config = map[string]int{"key": TestOne} // want "KTN-VAR-007"

	// Using variable to avoid unused warning
	_ = config
}

// badVarInIf shows incorrect use of var in if block.
// Local variables should use := even inside control structures.
func badVarInIf() {
	// Initialize test variable
	x := TestX

	// Check if x exceeds threshold
	if x > ThresholdValue {
		// Variable in if block (should use :=)
		var result = "big" // want "KTN-VAR-007"

		// Using variable to avoid unused warning
		_ = result
	}
}

// badVarInFor shows incorrect use of var in for loop.
// Local variables should use := even inside loops.
func badVarInFor() {
	// Loop through range
	for i := ZeroIndex; i < TestThree; i++ {
		// Variable in for loop (should use :=)
		var item = i * MultiplierValue // want "KTN-VAR-007"

		// Using variable to avoid unused warning
		_ = item
	}
}

// badVarInRange shows incorrect use of var in range loop.
// Local variables should use := even inside range loops.
func badVarInRange() {
	// Initialize slice for iteration
	nums := []int{TestOne, TestTwo, TestThree}

	// Iterate over slice
	for _, n := range nums {
		// Variable in range loop (should use :=)
		var doubled = n * MultiplierValue // want "KTN-VAR-007"

		// Using variable to avoid unused warning
		_ = doubled
	}
}

// badVarInSwitch shows incorrect use of var in switch.
// Local variables should use := even in switch cases.
func badVarInSwitch() {
	// Initialize test value
	x := ThresholdValue

	// Switch on value
	switch x {
	// Case when x equals threshold
	case ThresholdValue:
		// Variable in case block (should use :=)
		var msg = "five" // want "KTN-VAR-007"

		// Using variable to avoid unused warning
		_ = msg
	// Default case for other values
	default:
		// Variable in default block (should use :=)
		var other = "other" // want "KTN-VAR-007"

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
		var result = val * MultiplierValue // want "KTN-VAR-007"

		// Using variable to avoid unused warning
		_ = result
	// Default case when no channel operation is ready
	default:
		// Variable in default block (should use :=)
		var msg = "no data" // want "KTN-VAR-007"

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
