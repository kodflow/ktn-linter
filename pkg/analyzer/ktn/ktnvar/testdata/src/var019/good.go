// Package var017 provides good test cases.
package var017

import "sync"

const (
	// InitialValue is the initial counter value
	InitialValue int = 0
)

// init demonstrates correct usage of mutex pointers (compliant with KTN-VAR-019)
func init() {
	// Good: Using mutex pointer - no copy by value
	mu := &sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	// Good: Struct with mutex pointer field
	type counter struct {
		mu    *sync.Mutex // OK: Pointer to mutex
		value int
	}

	// Good: Initialize with mutex pointer
	c := &counter{
		mu:    &sync.Mutex{},
		value: InitialValue,
	}

	// Good: Use with pointer receiver (no copy)
	c.mu.Lock()
	c.value++
	c.mu.Unlock()

	// Good: RWMutex pointer
	rwMu := &sync.RWMutex{}
	rwMu.RLock()
	defer rwMu.RUnlock()
}
