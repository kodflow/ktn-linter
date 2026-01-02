// Package var036 contains test cases for KTN-VAR-036 rule.
package var036

const (
	// BadSampleValue represents a sample numeric value for testing.
	BadSampleValue int = 1
	// BadSecondValue represents a second sample value.
	BadSecondValue int = 2
	// BadThirdValue represents a third sample value.
	BadThirdValue int = 3
)

// badIndexOf searches for target in items and returns its index.
// This pattern should use slices.Index instead.
//
// Params:
//   - items: slice to search in.
//   - target: value to find.
//
// Returns:
//   - int: index of target or -1 if not found.
func badIndexOf(items []int, target int) int {
	// Iterate over all items with index
	for i, v := range items {
		// Check if value equals target
		if v == target { // want "KTN-VAR-036"
			// Return the index
			return i
		}
	}
	// Not found, return -1
	return -1
}

// badFindIndex searches for val in data and returns its index.
// This pattern should use slices.Index instead.
//
// Params:
//   - data: slice to search in.
//   - val: value to find.
//
// Returns:
//   - int: index of val or -1 if not found.
func badFindIndex(data []string, val string) int {
	// Iterate over all items with index
	for idx, item := range data {
		// Check if item equals val
		if item == val { // want "KTN-VAR-036"
			// Return the index
			return idx
		}
	}
	// Not found, return -1
	return -1
}

// init uses all the bad functions to avoid unused warnings.
func init() {
	// Create test slice with sample values
	ints := []int{BadSampleValue, BadSecondValue, BadThirdValue}
	// Test badIndexOf
	_ = badIndexOf(ints, BadSecondValue)

	// Create test string slice
	strs := []string{"a", "b", "c"}
	// Test badFindIndex
	_ = badFindIndex(strs, "b")
}
