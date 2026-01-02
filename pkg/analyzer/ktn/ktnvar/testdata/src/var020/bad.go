// Package var020 contains test cases for KTN-VAR-020.
package var020

// badExample demonstrates incorrect empty slice declarations.
// These patterns cause unnecessary allocations when a nil slice would suffice.
func badExample() {
	// Empty slice literal - prefer nil slice
	items := []string{} // want "KTN-VAR-020"
	// make with zero length and no capacity
	data := make([]int, 0) // want "KTN-VAR-020"
	// Another empty slice literal
	list := []int{} // want "KTN-VAR-020"

	// Use variables to avoid unused warning
	_ = items
	_ = data
	_ = list
}

// init calls the bad example function.
func init() {
	// Call the bad example function
	badExample()
}
