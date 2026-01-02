// Package var037 contains test cases for KTN-VAR-037 rule.
package var037

const (
	// BadSampleValueOne is a sample value for testing.
	BadSampleValueOne int = 1
	// BadSampleValueTwo is a sample value for testing.
	BadSampleValueTwo int = 2
	// BadSampleValueThree is a sample value for testing.
	BadSampleValueThree int = 3
)

// badGetKeys collects map keys using manual loop instead of maps.Keys().
//
// Params:
//   - m: map to get keys from.
//
// Returns:
//   - []string: slice of keys.
func badGetKeys(m map[string]int) []string {
	// Declare keys slice
	var keys []string
	// Iterate through map keys manually
	for k := range m {
		keys = append(keys, k) // want "KTN-VAR-037"
	}
	// Return collected keys
	return keys
}

// badGetValues collects map values using manual loop instead of maps.Values().
//
// Params:
//   - m: map to get values from.
//
// Returns:
//   - []int: slice of values.
func badGetValues(m map[string]int) []int {
	// Declare values slice
	var values []int
	// Iterate through map values manually
	for _, v := range m {
		values = append(values, v) // want "KTN-VAR-037"
	}
	// Return collected values
	return values
}

// init uses all the bad functions to avoid unused warnings.
func init() {
	// Create test map with sample values
	testMap := map[string]int{
		"a": BadSampleValueOne,
		"b": BadSampleValueTwo,
		"c": BadSampleValueThree,
	}

	// Test badGetKeys function
	_ = badGetKeys(testMap)

	// Test badGetValues function
	_ = badGetValues(testMap)
}
