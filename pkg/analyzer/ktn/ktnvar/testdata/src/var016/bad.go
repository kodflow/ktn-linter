// Bad examples for the var017 test case.
package var016

// Bad: Using make([]T, N) where N is small constant

const (
	// FixedSize est la taille fixe pour le test
	FixedSize int = 10
	// SmallSize est la petite taille pour le test
	SmallSize int = 64
	// TinySize est la très petite taille pour le test
	TinySize int = 5
	// MediumSize est la taille moyenne pour le test
	MediumSize int = 256
	// MaxSize est la taille maximale pour le test
	MaxSize int = 1024
)

// badFixedSize uses make with small constant size
func badFixedSize() {
	// Small fixed size should use array
	items := make([]int, FixedSize)
	_ = items
}

// badSmallBuffer creates slice with constant small size
func badSmallBuffer() {
	// Constant size 64 should be array
	buffer := make([]byte, SmallSize)
	_ = buffer
}

// badTinySlice creates very small slice
func badTinySlice() {
	// Very small size should use array
	data := make([]string, TinySize)
	_ = data
}

// badMediumSlice creates medium-sized slice
func badMediumSlice() {
	// Size 256 still small enough for stack
	values := make([]float64, MediumSize)
	_ = values
}

// badMaxAllowed creates slice at boundary
func badMaxAllowed() {
	// Size 1024 is at the limit
	large := make([]int32, MaxSize)
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
