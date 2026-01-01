// Package var018 contains test cases for KTN-VAR-018.
package var018

// Bad: Using make([]T, N) where total size <= 64 bytes

// badByteBuffer creates slice with 32 bytes (≤64)
func badByteBuffer() {
	// 32 bytes: should use [32]byte
	buf := make([]byte, 32) // want "KTN-VAR-018"
	_ = buf
}

// badIntArray creates slice with 8 ints (64 bytes on 64-bit)
func badIntArray() {
	// 8 ints * 8 bytes = 64 bytes
	arr := make([]int, 8) // want "KTN-VAR-018"
	_ = arr
}

// badSmallExact creates slice exactly 64 bytes
func badSmallExact() {
	// Exactly 64 bytes
	small := make([]byte, 64) // want "KTN-VAR-018"
	_ = small
}

// badTinySlice creates very small slice
func badTinySlice() {
	// 4 ints * 8 bytes = 32 bytes
	tiny := make([]int64, 4) // want "KTN-VAR-018"
	_ = tiny
}

// badSingleInt creates single element array
func badSingleInt() {
	// 1 int * 8 bytes = 8 bytes
	single := make([]int, 1) // want "KTN-VAR-018"
	_ = single
}
