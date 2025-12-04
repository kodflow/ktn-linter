// Bad examples for the var017 test case.
package var017

// Bad: Using make([]T, N) where N is small constant

const (
	// FIXED_SIZE est la taille fixe pour le test
	FIXED_SIZE int = 10
	// SMALL_SIZE est la petite taille pour le test
	SMALL_SIZE int = 64
	// TINY_SIZE est la très petite taille pour le test
	TINY_SIZE int = 5
	// MEDIUM_SIZE est la taille moyenne pour le test
	MEDIUM_SIZE int = 256
	// MAX_SIZE est la taille maximale pour le test
	MAX_SIZE int = 1024
)

// badFixedSize uses make with small constant size
func badFixedSize() {
	// Small fixed size should use array
	items := make([]int, FIXED_SIZE)
	_ = items
}

// badSmallBuffer creates slice with constant small size
func badSmallBuffer() {
	// Constant size 64 should be array
	buffer := make([]byte, SMALL_SIZE)
	_ = buffer
}

// badTinySlice creates very small slice
func badTinySlice() {
	// Very small size should use array
	data := make([]string, TINY_SIZE)
	_ = data
}

// badMediumSlice creates medium-sized slice
func badMediumSlice() {
	// Size 256 still small enough for stack
	values := make([]float64, MEDIUM_SIZE)
	_ = values
}

// badMaxAllowed creates slice at boundary
func badMaxAllowed() {
	// Size 1024 is at the limit
	large := make([]int32, MAX_SIZE)
	_ = large
}

// init utilise les fonctions privées
func init() {
	// Appel de badFixedSize
	badFixedSize()
	// Appel de badSmallBuffer
	badSmallBuffer()
	// Appel de badTinySlice
	badTinySlice()
	// Appel de badMediumSlice
	badMediumSlice()
	// Appel de badMaxAllowed
	badMaxAllowed()
}
