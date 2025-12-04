// Package func004_special tests special functions handling.
package func004_special

// deadFunction is never called - should trigger error. // want "KTN-FUNC-004"
func deadFunction() {
	// Dead code
}
