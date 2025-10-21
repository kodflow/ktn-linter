package var016

// Bad: Using make([]T, N) where N is small constant

// badFixedSize uses make with small constant size
func badFixedSize() {
	// Small fixed size should use array
	items := make([]int, 10) // want "KTN-VAR-016"
	_ = items
}

// badSmallBuffer creates slice with constant small size
func badSmallBuffer() {
	// Constant size 64 should be array
	buffer := make([]byte, 64) // want "KTN-VAR-016"
	_ = buffer
}

// badTinySlice creates very small slice
func badTinySlice() {
	// Very small size should use array
	data := make([]string, 5) // want "KTN-VAR-016"
	_ = data
}

// badMediumSlice creates medium-sized slice
func badMediumSlice() {
	// Size 256 still small enough for stack
	values := make([]float64, 256) // want "KTN-VAR-016"
	_ = values
}

// badMaxAllowed creates slice at boundary
func badMaxAllowed() {
	// Size 1024 is at the limit
	large := make([]int32, 1024) // want "KTN-VAR-016"
	_ = large
}
