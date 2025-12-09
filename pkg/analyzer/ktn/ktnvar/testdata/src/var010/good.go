// Good examples for the var010 test case.
package var010

import "sync"

const (
	// ValueThree is constant value 3
	ValueThree int = 3
	// ValueSixtyFour is constant value 64
	ValueSixtyFour int = 64
	// ValueHundred is constant value 100
	ValueHundred int = 100
	// Value1024 is constant value 1024
	Value1024 int = 1024
)

// bufferPool is a sync.Pool for byte buffers
var bufferPool *sync.Pool = &sync.Pool{
	New: func() any {
		// Buffer size optimized for common use case
		buffer := make([]byte, 0, Value1024)
		// Return preallocated buffer
		return buffer
	},
}

// init demonstrates correct usage patterns
func init() {
	// Using sync.Pool for buffer reuse
	for i := range ValueHundred {
		// Get buffer from pool
		buf := bufferPool.Get().([]byte)
		_ = buf
		// Put buffer back to pool
		bufferPool.Put(buf)
		_ = i
	}

	// Buffer allocated once before loop with array
	var buffer [Value1024]byte
	// Loop reuses buffer
	for i := range ValueHundred {
		_ = buffer
		_ = i
	}

	// Not in a loop, no pool needed, use array
	var bufferNoLoop [Value1024]byte
	_ = bufferNoLoop

	// Small fixed iteration count where pool overhead may not help with array
	var smallBuffer [ValueSixtyFour]byte
	// Loop with small iterations
	for i := range ValueThree {
		_ = smallBuffer
		_ = i
	}
}
