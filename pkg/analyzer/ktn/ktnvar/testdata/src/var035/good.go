// Package var035 contains test cases for KTN-VAR-035 rule.
package var035

import "slices"

const (
	// SampleOne represents a sample numeric value for testing.
	SampleOne int = 1
	// SampleTwo represents a second sample value.
	SampleTwo int = 2
	// SampleThree represents a third sample value.
	SampleThree int = 3
)

// goodContains checks if items slice contains target using slices.Contains.
//
// Params:
//   - items: slice to search in.
//   - target: value to search for.
//
// Returns:
//   - bool: true if target is found.
func goodContains(items []string, target string) bool {
	// Use slices.Contains for idiomatic search
	return slices.Contains(items, target)
}

// goodHasItem checks if data slice contains val using slices.Contains.
//
// Params:
//   - data: slice to search in.
//   - val: value to search for.
//
// Returns:
//   - bool: true if val is found.
func goodHasItem(data []int, val int) bool {
	// Use slices.Contains for idiomatic search
	return slices.Contains(data, val)
}

// goodComplexCheck performs a complex check that is not a simple contains.
//
// Params:
//   - items: slice to search in.
//   - target: value to compare against.
//
// Returns:
//   - bool: true if any item is greater than target.
func goodComplexCheck(items []int, target int) bool {
	// Iterate through all items in slice
	for _, v := range items {
		// Check if current item is greater than target (not equality)
		if v > target {
			// Found an item greater than target
			return true
		}
	}
	// No item greater than target found
	return false
}

// goodNotEqualCheck performs a not-equal check (not a contains pattern).
//
// Params:
//   - items: slice to search in.
//   - target: value to compare against.
//
// Returns:
//   - bool: true if any item is not equal to target.
func goodNotEqualCheck(items []int, target int) bool {
	// Iterate through all items in slice
	for _, v := range items {
		// Check if current item is not equal to target
		if v != target {
			// Found an item not equal to target
			return true
		}
	}
	// All items equal to target
	return false
}

// goodWithMultipleReturns has multiple return values (more complex than simple contains).
//
// Params:
//   - items: slice to search in.
//   - target: value to search for.
//
// Returns:
//   - bool: result of complex logic.
//   - int: the value found or zero.
func goodWithMultipleReturns(items []int, target int) (bool, int) {
	// Iterate through all items in slice
	for _, v := range items {
		// Check if current item equals target
		if v == target {
			// Found the target, return both values
			return true, v
		}
	}
	// Target not found, return zero value
	return false, 0
}

// goodUsingIndex uses the index in the comparison (not simple contains).
//
// Params:
//   - items: slice to search in.
//   - targetIdx: index to check.
//
// Returns:
//   - bool: true if item at index equals targetIdx.
func goodUsingIndex(items []int, targetIdx int) bool {
	// Iterate through all items with index
	for idx, v := range items {
		// Check if index equals target index and value is positive
		if idx == targetIdx && v > 0 {
			// Found matching index with positive value
			return true
		}
	}
	// No match found
	return false
}

// goodMultipleStatements has multiple statements in if body.
//
// Params:
//   - items: slice to search in.
//   - target: value to search for.
//
// Returns:
//   - bool: result of search.
func goodMultipleStatements(items []int, target int) bool {
	// Iterate through all items in slice
	for _, v := range items {
		// Check if current item equals target
		if v == target {
			// Print statement (multiple statements in body)
			_ = v
			// Found the target
			return true
		}
	}
	// Target not found
	return false
}

// init uses all the good functions to avoid unused warnings.
func init() {
	// Create test slice with sample values
	ints := []int{SampleOne, SampleTwo, SampleThree}

	// Test goodContains function
	_ = goodContains([]string{"a", "b", "c"}, "b")

	// Test goodHasItem function
	_ = goodHasItem(ints, SampleTwo)

	// Test goodComplexCheck function
	_ = goodComplexCheck(ints, SampleTwo)

	// Test goodNotEqualCheck function
	_ = goodNotEqualCheck(ints, SampleTwo)

	// Test goodWithMultipleReturns function
	_, _ = goodWithMultipleReturns(ints, SampleTwo)

	// Test goodUsingIndex function
	_ = goodUsingIndex(ints, SampleOne)

	// Test goodMultipleStatements function
	_ = goodMultipleStatements(ints, SampleTwo)
}
