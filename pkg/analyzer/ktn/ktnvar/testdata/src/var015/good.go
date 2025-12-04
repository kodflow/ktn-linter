// Good examples for the var015 test case.
package var015

import "sync"

const (
	// VALUE_THREE is constant value 3
	VALUE_THREE int = 3
	// VALUE_SIXTY_FOUR is constant value 64
	VALUE_SIXTY_FOUR int = 64
	// VALUE_HUNDRED is constant value 100
	VALUE_HUNDRED int = 100
	// VALUE_1024 is constant value 1024
	VALUE_1024 int = 1024
)

// Good: Using sync.Pool or creating buffers outside loops

// bufferPool is a sync.Pool for byte buffers
var bufferPool = &sync.Pool{
	New: func() any {
		// Buffer size optimized for common use case
		buffer := make([]byte, 0, VALUE_1024)
		// Return preallocated buffer
		return buffer
	},
}

// goodWithPool uses sync.Pool for buffer reuse
func goodWithPool() {
	// Loop processes items
	for i := range VALUE_HUNDRED {
		// Get buffer from pool
		buffer := bufferPool.Get().([]byte)
		_ = buffer
		// Put buffer back to pool
		bufferPool.Put(buffer)
		// Utilisation de i pour éviter le warning
		_ = i
	}
}

// goodOutsideLoop creates buffer outside the loop
func goodOutsideLoop() {
	// Buffer allocated once before loop with array
	var buffer [VALUE_1024]byte
	// Loop reuses buffer
	for i := range VALUE_HUNDRED {
		_ = buffer
		// Utilisation de i pour éviter le warning
		_ = i
	}
}

// goodNoLoop creates buffer outside loop context
func goodNoLoop() {
	// Not in a loop, no pool needed, use array
	var buffer [VALUE_1024]byte
	_ = buffer
}

// goodSmallLoop creates buffer where pooling overhead not justified
func goodSmallLoop() {
	// Small fixed iteration count where pool overhead may not help with array
	var buffer [VALUE_SIXTY_FOUR]byte
	// Very small fixed loop
	for i := range VALUE_THREE {
		_ = buffer
		_ = i
	}
}

// init utilise les fonctions privées
func init() {
	// Appel de goodWithPool
	goodWithPool()
	// Appel de goodOutsideLoop
	goodOutsideLoop()
	// Appel de goodNoLoop
	goodNoLoop()
	// Appel de goodSmallLoop
	goodSmallLoop()
}
