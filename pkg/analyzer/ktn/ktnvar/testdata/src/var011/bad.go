// Bad examples for the var011 test case.
package var011

// Bad: Creating []byte buffers repeatedly in loops without sync.Pool

const (
	LOOP_ITERATIONS         int = 100
	BUFFER_SIZE             int = 1024
	SMALL_BUFFER_SIZE       int = 512
	TINY_BUFFER_SIZE        int = 128
	MEDIUM_BUFFER_SIZE      int = 256
	CONST_BUFFER_SIZE       int = 2048
	OUTER_ITERATIONS        int = 10
	INNER_ITERATIONS        int = 10
	NESTED_BUFFER_SIZE      int = 64
	LARGE_BUFFER_SIZE       int = 512
	LARGE_BUFFER_CAPACITY   int = 1024
	ITEMS_COUNT             int = 3
	COUNTER_LIMIT           int = 100
	RANGE_ITERATIONS        int = 50
)

// badProcessInLoop creates buffer in loop without sync.Pool
func badProcessInLoop() {
	items := make([]int, 0, LOOP_ITERATIONS)
	// Loop processes items
	for range LOOP_ITERATIONS {
		// Buffer created repeatedly in loop
		buffer := make([]byte, 0, BUFFER_SIZE)
		items = append(items, len(buffer))
	}
}

// badRangeLoop creates buffer in range loop without sync.Pool
func badRangeLoop() {
	items := [ITEMS_COUNT]string{"a", "b", "c"}
	results := make([]int, 0, ITEMS_COUNT)
	// Loop processes items
	for _, item := range items {
		// Buffer created in each iteration
		buf := make([]byte, 0, SMALL_BUFFER_SIZE)
		_ = item
		results = append(results, len(buf))
	}
}

// badInfiniteLoop creates buffer in infinite loop
func badInfiniteLoop() {
	// Infinite loop processing
	for {
		// Buffer allocated every iteration
		b := make([]byte, 0, TINY_BUFFER_SIZE)
		_ = b
		break
	}
}

// badWhileStyle creates buffer in while-style loop
func badWhileStyle() {
	counter := 0
	results := make([]int, 0, COUNTER_LIMIT)
	// While-style loop
	for counter < COUNTER_LIMIT {
		// Buffer created each iteration
		data := make([]byte, 0, MEDIUM_BUFFER_SIZE)
		results = append(results, len(data))
		counter++
	}
}

// badConstSizeBuffer creates buffer with const size in loop
func badConstSizeBuffer() {
	results := make([]int, 0, RANGE_ITERATIONS)
	// Loop processes items
	for range RANGE_ITERATIONS {
		// Buffer with const size
		buf := make([]byte, 0, CONST_BUFFER_SIZE)
		results = append(results, len(buf))
	}
}

// badNestedLoop creates buffer in nested loop
func badNestedLoop() {
	// Outer loop
	for range OUTER_ITERATIONS {
		// Inner loop
		for range INNER_ITERATIONS {
			// Buffer allocated in nested loop
			temp := make([]byte, 0, NESTED_BUFFER_SIZE)
			_ = temp
		}
	}
}

// badMakeWithCapacity creates buffer with capacity in loop
func badMakeWithCapacity() {
	results := make([]int, 0, LOOP_ITERATIONS)
	// Loop processes items
	for range LOOP_ITERATIONS {
		// Buffer with both length and capacity
		buffer := make([]byte, LARGE_BUFFER_SIZE, LARGE_BUFFER_CAPACITY)
		results = append(results, len(buffer))
	}
}

// init utilise les fonctions privées
func init() {
	// Appel de badProcessInLoop
	badProcessInLoop()
	// Appel de badRangeLoop
	badRangeLoop()
	// Appel de badInfiniteLoop
	badInfiniteLoop()
	// Appel de badWhileStyle
	badWhileStyle()
	// Appel de badConstSizeBuffer
	badConstSizeBuffer()
	// Appel de badNestedLoop
	badNestedLoop()
	// Appel de badMakeWithCapacity
	badMakeWithCapacity()
}
