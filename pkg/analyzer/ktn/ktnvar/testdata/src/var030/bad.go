// Package var030 contains test cases for KTN-VAR-030 rule.
package var030

const (
	// BadSampleValue represents a sample numeric value for testing.
	BadSampleValue int = 1
	// BadSecondValue represents a second sample value.
	BadSecondValue int = 2
	// BadThirdValue represents a third sample value.
	BadThirdValue int = 3
)

// badMakeCopyInt copies a slice using make+copy pattern (should use slices.Clone).
//
// Params:
//   - original: slice to copy.
//
// Returns:
//   - []int: copied slice.
func badMakeCopyInt(original []int) []int {
	// Create a new slice with the same length
	clone := make([]int, len(original))
	// Copy elements from original to clone
	copy(clone, original) // want "KTN-VAR-030"
	// Return the cloned slice
	return clone
}

// badMakeCopyString copies a string slice using make+copy pattern.
//
// Params:
//   - original: slice to copy.
//
// Returns:
//   - []string: copied slice.
func badMakeCopyString(original []string) []string {
	// Create a new slice with the same length
	clone := make([]string, len(original))
	// Copy elements from original to clone
	copy(clone, original) // want "KTN-VAR-030"
	// Return the cloned slice
	return clone
}

// badAppendNilInt clones a slice using append([]T(nil), s...) pattern.
//
// Params:
//   - original: slice to copy.
//
// Returns:
//   - []int: copied slice.
func badAppendNilInt(original []int) []int {
	// Append nil pattern for cloning
	return append([]int(nil), original...) // want "KTN-VAR-030"
}

// badAppendNilString clones a string slice using append([]T(nil), s...) pattern.
//
// Params:
//   - original: slice to copy.
//
// Returns:
//   - []string: copied slice.
func badAppendNilString(original []string) []string {
	// Append nil pattern for cloning
	return append([]string(nil), original...) // want "KTN-VAR-030"
}

// badMakeCopyBytes copies a byte slice using make+copy pattern.
//
// Params:
//   - original: slice to copy.
//
// Returns:
//   - []byte: copied slice.
func badMakeCopyBytes(original []byte) []byte {
	// Create a new slice with the same length
	result := make([]byte, len(original))
	// Copy elements from original to result
	copy(result, original) // want "KTN-VAR-030"
	// Return the cloned slice
	return result
}

// badAppendNilBytes clones a byte slice using append([]T(nil), s...) pattern.
//
// Params:
//   - original: slice to copy.
//
// Returns:
//   - []byte: copied slice.
func badAppendNilBytes(original []byte) []byte {
	// Append nil pattern for cloning
	return append([]byte(nil), original...) // want "KTN-VAR-030"
}

// init uses all the bad functions to avoid unused warnings.
func init() {
	// Create test slice with sample values
	ints := []int{BadSampleValue, BadSecondValue, BadThirdValue}
	// Test badMakeCopyInt
	_ = badMakeCopyInt(ints)

	// Create test string slice
	strs := []string{"a", "b", "c"}
	// Test badMakeCopyString
	_ = badMakeCopyString(strs)

	// Test badAppendNilInt
	_ = badAppendNilInt(ints)

	// Test badAppendNilString
	_ = badAppendNilString(strs)

	// Create test byte slice
	bytes := []byte{byte(BadSampleValue), byte(BadSecondValue), byte(BadThirdValue)}
	// Test badMakeCopyBytes
	_ = badMakeCopyBytes(bytes)

	// Test badAppendNilBytes
	_ = badAppendNilBytes(bytes)
}
