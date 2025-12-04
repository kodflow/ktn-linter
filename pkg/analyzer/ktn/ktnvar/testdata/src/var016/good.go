// Good examples for the var017 test case.
package var016

// Good: Using arrays for small fixed sizes or slices when appropriate

const (
	// ARRAY_SIZE_SMALL is small array size
	ARRAY_SIZE_SMALL int = 10
	// BUFFER_SIZE is buffer size
	BUFFER_SIZE int = 64
	// LARGE_SIZE is large allocation size
	LARGE_SIZE int = 2048
	// CAPACITY_MEDIUM is medium capacity
	CAPACITY_MEDIUM int = 100
	// LENGTH_SMALL is small length
	LENGTH_SMALL int = 10
	// CAPACITY_SMALL is small capacity
	CAPACITY_SMALL int = 20
)

// goodArray uses array for fixed small size
func goodArray() {
	// Array allocated on stack
	var items [ARRAY_SIZE_SMALL]int
	_ = items
}

// goodArrayBuffer uses array for buffer
func goodArrayBuffer() {
	// Fixed size buffer as array
	var buffer [BUFFER_SIZE]byte
	_ = buffer
}

// goodDynamicSize uses make with variable size
//
// Params:
//   - n: dynamic size for slice allocation
func goodDynamicSize(n int) {
	// Size is variable, must use slice with capacity
	items := make([]int, 0, n)
	items = append(items, n)
	_ = items
}

// goodLargeSlice uses make for large allocation
func goodLargeSlice() {
	// Size > 1024, slice is appropriate
	large := make([]byte, 0, LARGE_SIZE)
	_ = large
}

// goodGrowingSlice uses make with capacity for growing slice
func goodGrowingSlice() {
	// Slice will grow, needs heap allocation
	items := make([]string, 0, CAPACITY_MEDIUM)
	items = append(items, "test")
	_ = items
}

// goodMakeWithCapacity uses make with different len/cap
func goodMakeWithCapacity() {
	// Different length and capacity
	items := make([]int, 0, CAPACITY_SMALL)
	items = append(items, LENGTH_SMALL)
	_ = items
}

// init utilise les fonctions privées
func init() {
	// Appel de goodArray
	goodArray()
	// Appel de goodArrayBuffer
	goodArrayBuffer()
	// Appel de goodDynamicSize
	goodDynamicSize(0)
	// Appel de goodLargeSlice
	goodLargeSlice()
	// Appel de goodGrowingSlice
	goodGrowingSlice()
	// Appel de goodMakeWithCapacity
	goodMakeWithCapacity()
}
