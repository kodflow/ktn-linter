// Package var036 contains test cases for KTN-VAR-036 rule.
package var036

import "slices"

const (
	// SampleValue represents a sample numeric value for testing.
	SampleValue int = 1
	// SecondValue represents a second sample value.
	SecondValue int = 2
	// ThirdValue represents a third sample value.
	ThirdValue int = 3
)

// goodIndexOf uses slices.Index to find target in items.
//
// Params:
//   - items: slice to search in.
//   - target: value to find.
//
// Returns:
//   - int: index of target or -1 if not found.
func goodIndexOf(items []int, target int) int {
	// Use slices.Index for idiomatic search
	return slices.Index(items, target)
}

// goodFindIndex uses slices.Index to find val in data.
//
// Params:
//   - data: slice to search in.
//   - val: value to find.
//
// Returns:
//   - int: index of val or -1 if not found.
func goodFindIndex(data []string, val string) int {
	// Use slices.Index for idiomatic search
	return slices.Index(data, val)
}

// goodComplexFind uses a custom comparison (not simple equality).
// This should NOT trigger the rule since it uses > instead of ==.
//
// Params:
//   - items: slice to search in.
//   - threshold: value to compare against.
//
// Returns:
//   - int: index of first element greater than threshold, or -1.
func goodComplexFind(items []int, threshold int) int {
	// Iterate over all items with index
	for i, v := range items {
		// Check if value is greater than threshold (not simple equality)
		if v > threshold {
			// Return the index
			return i
		}
	}
	// Not found, return -1
	return -1
}

// goodRangeWithoutValue iterates without value variable.
// This should NOT trigger the rule since there's no value comparison.
//
// Params:
//   - items: slice to iterate over.
//
// Returns:
//   - int: always returns -1 for this example.
func goodRangeWithoutValue(items []int) int {
	// Iterate over indices only
	for i := range items {
		// Use index directly
		_ = i
	}
	// Return -1
	return -1
}

// goodNotReturningMinusOne has a different return value.
// This should NOT trigger the rule since it returns 0 instead of -1.
//
// Params:
//   - items: slice to search in.
//   - target: value to find.
//
// Returns:
//   - int: index of target or 0 if not found.
func goodNotReturningMinusOne(items []int, target int) int {
	// Iterate over all items with index
	for i, v := range items {
		// Check if value equals target
		if v == target {
			// Return the index
			return i
		}
	}
	// Different return value - not the -1 pattern
	return 0
}

// init uses all the good functions to avoid unused warnings.
func init() {
	// Create test slice with sample values
	ints := []int{SampleValue, SecondValue, ThirdValue}
	// Test goodIndexOf
	_ = goodIndexOf(ints, SecondValue)

	// Create test string slice
	strs := []string{"a", "b", "c"}
	// Test goodFindIndex
	_ = goodFindIndex(strs, "b")

	// Test goodComplexFind
	_ = goodComplexFind(ints, SecondValue)

	// Test goodRangeWithoutValue
	_ = goodRangeWithoutValue(ints)

	// Test goodNotReturningMinusOne
	_ = goodNotReturningMinusOne(ints, SecondValue)
}
