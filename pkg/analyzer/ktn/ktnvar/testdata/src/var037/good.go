// Package var037 contains test cases for KTN-VAR-037 rule.
package var037

import (
	"maps"
	"slices"
)

const (
	// SuffixValue is a sample suffix used for testing transformations.
	SuffixValue string = "_suffix"
	// SampleValueOne is a sample value for testing.
	SampleValueOne int = 1
	// SampleValueTwo is a sample value for testing.
	SampleValueTwo int = 2
	// SampleValueThree is a sample value for testing.
	SampleValueThree int = 3
	// MultiplierFactor is the multiplier used for doubling values.
	MultiplierFactor int = 2
)

// goodGetKeys collects map keys using slices.Collect and maps.Keys.
//
// Params:
//   - m: map to get keys from.
//
// Returns:
//   - []string: slice of keys.
func goodGetKeys(m map[string]int) []string {
	// Use idiomatic Go 1.23+ pattern
	return slices.Collect(maps.Keys(m))
}

// goodGetValues collects map values using slices.Collect and maps.Values.
//
// Params:
//   - m: map to get values from.
//
// Returns:
//   - []int: slice of values.
func goodGetValues(m map[string]int) []int {
	// Use idiomatic Go 1.23+ pattern
	return slices.Collect(maps.Values(m))
}

// goodTransformKeys transforms keys with suffix (not simple collection).
//
// Params:
//   - m: map to get keys from.
//
// Returns:
//   - []string: slice of transformed keys.
func goodTransformKeys(m map[string]int) []string {
	// Declare keys slice
	var keys []string
	// Iterate through map keys with transformation
	for k := range m {
		// Add suffix to key (transformation, not simple collection)
		keys = append(keys, k+SuffixValue)
	}
	// Return transformed keys
	return keys
}

// goodTransformValues transforms values by doubling (not simple collection).
//
// Params:
//   - m: map to get values from.
//
// Returns:
//   - []int: slice of transformed values.
func goodTransformValues(m map[string]int) []int {
	// Declare values slice
	var values []int
	// Iterate through map values with transformation
	for _, v := range m {
		// Double the value (transformation, not simple collection)
		values = append(values, v*MultiplierFactor)
	}
	// Return transformed values
	return values
}

// goodWithCondition filters values (not simple collection).
//
// Params:
//   - m: map to filter values from.
//
// Returns:
//   - []int: slice of filtered values.
func goodWithCondition(m map[string]int) []int {
	// Declare values slice
	var values []int
	// Iterate through map with filtering
	for _, v := range m {
		// Only add positive values
		if v > 0 {
			values = append(values, v)
		}
	}
	// Return filtered values
	return values
}

// goodUsingBothKeyValue uses both key and value (not simple collection).
//
// Params:
//   - m: map to process.
//
// Returns:
//   - []string: slice of formatted entries.
func goodUsingBothKeyValue(m map[string]int) []string {
	// Declare results slice
	var results []string
	// Iterate through map using both key and value
	for k, v := range m {
		// Both key and value are used
		if v > 0 {
			results = append(results, k)
		}
	}
	// Return results
	return results
}

// goodMultipleStatements has multiple statements in loop body.
//
// Params:
//   - m: map to get keys from.
//
// Returns:
//   - []string: slice of keys.
func goodMultipleStatements(m map[string]int) []string {
	// Declare keys slice
	var keys []string
	// Declare count variable
	count := 0
	// Iterate through map with multiple statements
	for k := range m {
		// Increment count
		count++
		// Add key to slice
		keys = append(keys, k)
	}
	// Use count to avoid unused warning
	_ = count
	// Return keys
	return keys
}

// goodAppendMultiple appends multiple elements (not simple pattern).
//
// Params:
//   - m: map to get keys from.
//
// Returns:
//   - []string: slice with duplicated keys.
func goodAppendMultiple(m map[string]int) []string {
	// Declare keys slice
	var keys []string
	// Iterate through map
	for k := range m {
		// Append key twice (variadic append)
		keys = append(keys, k, k)
	}
	// Return keys
	return keys
}

// init uses all the good functions to avoid unused warnings.
func init() {
	// Create test map with sample values
	testMap := map[string]int{
		"a": SampleValueOne,
		"b": SampleValueTwo,
		"c": SampleValueThree,
	}

	// Test goodGetKeys function
	_ = goodGetKeys(testMap)

	// Test goodGetValues function
	_ = goodGetValues(testMap)

	// Test goodTransformKeys function
	_ = goodTransformKeys(testMap)

	// Test goodTransformValues function
	_ = goodTransformValues(testMap)

	// Test goodWithCondition function
	_ = goodWithCondition(testMap)

	// Test goodUsingBothKeyValue function
	_ = goodUsingBothKeyValue(testMap)

	// Test goodMultipleStatements function
	_ = goodMultipleStatements(testMap)

	// Test goodAppendMultiple function
	_ = goodAppendMultiple(testMap)
}
