// Package var029 contains test cases for KTN-VAR-029 rule.
package var029

const (
	// BadSampleValue represents a sample numeric value for testing.
	BadSampleValue int = 1
	// BadSecondValue represents a second sample value.
	BadSecondValue int = 2
	// BadThirdValue represents a third sample value.
	BadThirdValue int = 3
)

// badGrow grows a slice manually instead of using slices.Grow.
//
// Params:
//   - s: slice to grow.
//   - n: additional capacity needed.
//
// Returns:
//   - []int: grown slice.
func badGrow(s []int, n int) []int {
	// Check if capacity needs to grow
	if cap(s)-len(s) < n { // want "KTN-VAR-029"
		// Create a new slice with increased capacity
		newSlice := make([]int, len(s), len(s)+n)
		// Copy existing elements to the new slice
		copy(newSlice, s)
		// Reassign to the original variable
		s = newSlice
	}
	// Return the potentially grown slice
	return s
}

// badGrowString grows a string slice manually.
//
// Params:
//   - s: slice to grow.
//   - n: additional capacity needed.
//
// Returns:
//   - []string: grown slice.
func badGrowString(s []string, n int) []string {
	// Check if capacity needs to grow
	if cap(s)-len(s) < n { // want "KTN-VAR-029"
		// Create a new slice with increased capacity
		newSlice := make([]string, len(s), len(s)+n)
		// Copy existing elements to the new slice
		copy(newSlice, s)
		// Reassign to the original variable
		s = newSlice
	}
	// Return the potentially grown slice
	return s
}

// badGrowBytes grows a byte slice manually.
//
// Params:
//   - s: slice to grow.
//   - n: additional capacity needed.
//
// Returns:
//   - []byte: grown slice.
func badGrowBytes(s []byte, n int) []byte {
	// Check if capacity needs to grow
	if cap(s)-len(s) < n { // want "KTN-VAR-029"
		// Create a new slice with increased capacity
		newSlice := make([]byte, len(s), len(s)+n)
		// Copy existing elements to the new slice
		copy(newSlice, s)
		// Reassign to the original variable
		s = newSlice
	}
	// Return the potentially grown slice
	return s
}

// init uses all the bad functions to avoid unused warnings.
func init() {
	// Create test slice with sample values
	ints := []int{BadSampleValue, BadSecondValue, BadThirdValue}
	// Test badGrow
	_ = badGrow(ints, BadSampleValue)

	// Create test string slice
	strs := []string{"a", "b", "c"}
	// Test badGrowString
	_ = badGrowString(strs, BadSampleValue)

	// Create test byte slice
	bytes := []byte{byte(BadSampleValue), byte(BadSecondValue)}
	// Test badGrowBytes
	_ = badGrowBytes(bytes, BadSampleValue)
}
