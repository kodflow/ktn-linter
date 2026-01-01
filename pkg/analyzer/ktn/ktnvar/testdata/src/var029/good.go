// Package var029 contains test cases for KTN-VAR-029 rule.
package var029

import "slices"

const (
	// SampleValue represents a sample numeric value for testing.
	SampleValue int = 1
	// SecondValue represents a second sample value.
	SecondValue int = 2
	// ThirdValue represents a third sample value.
	ThirdValue int = 3
)

// goodGrow grows a slice using slices.Grow (Go 1.21+).
//
// Params:
//   - s: slice to grow.
//   - n: additional capacity needed.
//
// Returns:
//   - []int: grown slice.
func goodGrow(s []int, n int) []int {
	// Use slices.Grow for idiomatic slice growing
	return slices.Grow(s, n)
}

// goodAppend appends items to a slice (not a grow pattern).
//
// Params:
//   - s: slice to append to.
//   - items: items to append.
//
// Returns:
//   - []int: slice with appended items.
func goodAppend(s []int, items []int) []int {
	// Simple append is OK
	return append(s, items...)
}

// goodMakeWithCapacity creates a slice with capacity (not a grow pattern).
//
// Params:
//   - n: capacity to allocate.
//
// Returns:
//   - []int: new slice with capacity.
func goodMakeWithCapacity(n int) []int {
	// Make with capacity is OK when not in grow pattern
	return make([]int, 0, n)
}

// goodDifferentCondition uses a different condition (not our pattern).
//
// Params:
//   - s: slice to check.
//   - n: threshold value.
//
// Returns:
//   - []int: processed slice.
func goodDifferentCondition(s []int, n int) []int {
	// Different condition: len < n (not cap-len < n)
	if len(s) < n {
		// Just return the original slice
		return s
	}
	// Return the slice unchanged
	return s
}

// goodSimpleCopy copies a slice using slices.Clone.
//
// Params:
//   - s: source slice.
//
// Returns:
//   - []int: copied slice.
func goodSimpleCopy(s []int) []int {
	// Use slices.Clone for copying
	return slices.Clone(s)
}

// init uses all the good functions to avoid unused warnings.
func init() {
	// Create test slice with sample values
	ints := []int{SampleValue, SecondValue, ThirdValue}
	// Test goodGrow
	_ = goodGrow(ints, SampleValue)

	// Create items to append
	items := []int{SampleValue, SecondValue}
	// Test goodAppend
	_ = goodAppend(ints, items)

	// Test goodMakeWithCapacity
	_ = goodMakeWithCapacity(ThirdValue)

	// Test goodDifferentCondition
	_ = goodDifferentCondition(ints, ThirdValue)

	// Test goodSimpleCopy
	_ = goodSimpleCopy(ints)
}
