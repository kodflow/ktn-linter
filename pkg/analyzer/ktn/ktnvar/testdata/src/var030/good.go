// Package var030 contains test cases for KTN-VAR-030 rule.
package var030

import "slices"

const (
	// SampleValue represents a sample numeric value for testing.
	SampleValue int = 1
	// SecondValue represents a second sample value.
	SecondValue int = 2
	// ThirdValue represents a third sample value.
	ThirdValue int = 3
)

// goodCloneInt clones a slice using slices.Clone (Go 1.21+).
//
// Params:
//   - original: slice to copy.
//
// Returns:
//   - []int: copied slice.
func goodCloneInt(original []int) []int {
	// Use slices.Clone for idiomatic cloning
	return slices.Clone(original)
}

// goodCloneString clones a string slice using slices.Clone.
//
// Params:
//   - original: slice to copy.
//
// Returns:
//   - []string: copied slice.
func goodCloneString(original []string) []string {
	// Use slices.Clone for idiomatic cloning
	return slices.Clone(original)
}

// goodCloneBytes clones a byte slice using slices.Clone.
//
// Params:
//   - original: slice to copy.
//
// Returns:
//   - []byte: copied slice.
func goodCloneBytes(original []byte) []byte {
	// Use slices.Clone for idiomatic cloning
	return slices.Clone(original)
}

// goodAppendWithElements appends to non-nil slice (not a clone pattern).
//
// Params:
//   - original: source slice.
//
// Returns:
//   - []int: new slice with additional elements.
func goodAppendWithElements(original []int) []int {
	// Start with a non-empty slice containing zero
	base := []int{0}
	// Append original elements to base
	return append(base, original...)
}

// goodMakeWithCapacity uses make with capacity (not a clone pattern).
// It builds a new slice by appending each element.
//
// Params:
//   - original: source slice.
//
// Returns:
//   - []int: new slice with same elements.
func goodMakeWithCapacity(original []int) []int {
	// Make with zero length and capacity equal to original length
	result := make([]int, 0, len(original))
	// Append elements one by one from original
	for _, v := range original {
		// Append each value to result
		result = append(result, v)
	}
	// Return the built slice
	return result
}

// init uses all the good functions to avoid unused warnings.
func init() {
	// Create test slice with sample values
	ints := []int{SampleValue, SecondValue, ThirdValue}
	// Test goodCloneInt
	_ = goodCloneInt(ints)

	// Create test string slice
	strs := []string{"a", "b", "c"}
	// Test goodCloneString
	_ = goodCloneString(strs)

	// Create test byte slice with sample values
	bytes := []byte{byte(SampleValue), byte(SecondValue), byte(ThirdValue)}
	// Test goodCloneBytes
	_ = goodCloneBytes(bytes)

	// Test goodAppendWithElements
	_ = goodAppendWithElements(ints)

	// Test goodMakeWithCapacity
	_ = goodMakeWithCapacity(ints)
}
