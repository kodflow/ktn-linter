// Package var020 contains test cases for KTN-VAR-020.
package var020

const (
	// DefaultCapacity is the default slice capacity
	DefaultCapacity int = 10
)

// goodExample demonstrates correct slice declarations.
// Nil slices and slices with capacity are preferred.
func goodExample() {
	// OK - nil slice (preferred)
	var items []string
	// OK - has capacity specified
	data := make([]int, 0, DefaultCapacity)
	// OK - nil slice (preferred)
	var list []int
	// OK - slice with elements
	initial := []string{"a", "b", "c"}

	// Use variables to avoid unused warning
	_ = items
	_ = data
	_ = list
	_ = initial
}

// init calls the good example function.
func init() {
	// Call the good example function
	goodExample()
}
