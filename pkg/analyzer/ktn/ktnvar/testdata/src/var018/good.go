// Package var018 provides good test cases.
package var018

// Good: Using arrays for ≤64 bytes or slices when appropriate

// goodLargeBuffer creates slice >64 bytes
func goodLargeBuffer() {
	// 128 bytes: heap acceptable (>64 bytes)
	buf := make([]byte, 128)
	_ = buf
}

// goodAlreadyArray uses array correctly
func goodAlreadyArray() {
	// Already using array syntax
	var small [32]byte
	_ = small
}

// goodDynamicSize uses dynamic size
func goodDynamicSize() {
	n := 32
	// Dynamic size: must use slice
	dynamic := make([]byte, n)
	_ = dynamic
}

// goodLargeIntArray creates slice >64 bytes
func goodLargeIntArray() {
	// 100 ints * 8 bytes = 800 bytes (>64 bytes)
	large := make([]int, 100)
	_ = large
}

// goodWithCapacity uses different length and capacity
func goodWithCapacity() {
	// Different capacity: needs slice semantics
	withCap := make([]byte, 32, 64)
	_ = withCap
}

// goodGrowingSlice creates slice that will grow
func goodGrowingSlice() {
	// Slice will grow: needs heap allocation
	growing := make([]string, 0, 10)
	growing = append(growing, "test")
	_ = growing
}
