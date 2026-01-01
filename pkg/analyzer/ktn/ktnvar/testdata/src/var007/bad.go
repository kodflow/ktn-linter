package var007

import (
	"bytes"
	"sync"
)

// badExample shows incorrect use of var with explicit initialization.
// Should use := instead of var when initializing variables.
func badExample() {
	// These should use := instead
	var x int = 42        // want "KTN-VAR-007"
	var s string = "test" // want "KTN-VAR-007"
	var err error = nil   // want "KTN-VAR-007"

	// Using variables to avoid unused warnings
	_ = x
	_ = s
	_ = err
}

// init ensures functions are used to avoid unused warnings
func init() {
	// Call badExample
	badExample()

	// Keep unused imports
	_ = bytes.Buffer{}
	_ = sync.WaitGroup{}
}
