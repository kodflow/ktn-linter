// Package var035 contains test cases for KTN-VAR-035 rule.
package var035

const (
	// BadSampleOne represents a sample numeric value for testing.
	BadSampleOne int = 1
	// BadSampleTwo represents a second sample value.
	BadSampleTwo int = 2
	// BadSampleThree represents a third sample value.
	BadSampleThree int = 3
)

// badContains checks if items slice contains target using manual pattern.
//
// Params:
//   - items: slice to search in.
//   - target: value to search for.
//
// Returns:
//   - bool: true if target is found.
func badContains(items []string, target string) bool {
	// Iterate through all items in slice
	for _, v := range items {
		// Check if current item equals target
		if v == target { // want "KTN-VAR-035"
			// Found the target
			return true
		}
	}
	// Target not found in slice
	return false
}

// badHasItem checks if data slice contains val using manual pattern.
//
// Params:
//   - data: slice to search in.
//   - val: value to search for.
//
// Returns:
//   - bool: true if val is found.
func badHasItem(data []int, val int) bool {
	// Iterate through all items in slice
	for _, item := range data {
		// Check if current item equals val
		if item == val { // want "KTN-VAR-035"
			// Found the value
			return true
		}
	}
	// Value not found in slice
	return false
}

// init uses all the bad functions to avoid unused warnings.
func init() {
	// Test badContains function
	_ = badContains([]string{"a", "b", "c"}, "b")

	// Test badHasItem function
	_ = badHasItem([]int{BadSampleOne, BadSampleTwo, BadSampleThree}, BadSampleTwo)
}
