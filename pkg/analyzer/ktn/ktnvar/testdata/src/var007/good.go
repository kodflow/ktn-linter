package var007

import (
	"bytes"
	"sync"
)

// goodExample demonstrates correct usage of variable declarations.
// Zero value declarations with var are idiomatic and acceptable.
func goodExample() {
	// OK - zero value intentional (no initialization)
	var err error
	// OK - zero value usable
	var buf bytes.Buffer
	// OK - zero value
	var wg sync.WaitGroup
	// OK - short syntax for initialization
	x := 42

	// Using variables to avoid unused warnings
	_ = err
	_ = buf
	_ = wg
	_ = x
}

// init ensures functions are used to avoid unused warnings
func init() {
	// Call goodExample
	goodExample()
}
