// Good examples for the var017 test case.
package var016

// Good: Using arrays for small fixed sizes or slices when appropriate

const (
	// ArraySizeSmall is small array size
	ArraySizeSmall int = 10
	// BufferSize is buffer size
	BufferSize int = 64
	// LargeSize is large allocation size
	LargeSize int = 2048
	// CapacityMedium is medium capacity
	CapacityMedium int = 100
	// LengthSmall is small length
	LengthSmall int = 10
	// CapacitySmall is small capacity
	CapacitySmall int = 20
)

// init demonstrates good practices for array and slice allocation
func init() {
	// Array allocated on stack - good for small fixed sizes
	var items [ArraySizeSmall]int
	_ = items

	// Fixed size buffer as array
	var buffer [BufferSize]byte
	_ = buffer

	// Dynamic size - use slice with capacity
	dynamicItems := make([]int, 0, CapacitySmall)
	dynamicItems = append(dynamicItems, LengthSmall)
	_ = dynamicItems

	// Size > 1024, slice is appropriate
	large := make([]byte, 0, LargeSize)
	_ = large

	// Slice will grow, needs heap allocation
	growingItems := make([]string, 0, CapacityMedium)
	growingItems = append(growingItems, "test")
	_ = growingItems

	// Different length and capacity
	withCapacity := make([]int, 0, CapacitySmall)
	withCapacity = append(withCapacity, LengthSmall)
	_ = withCapacity
}
