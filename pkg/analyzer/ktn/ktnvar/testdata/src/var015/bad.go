package var015

// Bad: Creating []byte buffers repeatedly in loops without sync.Pool

// badProcessInLoop creates buffer in loop without sync.Pool
func badProcessInLoop() {
	// Loop processes items
	for i := 0; i < 100; i++ {
		// Buffer created repeatedly in loop
		buffer := make([]byte, 1024)
		_ = buffer
	}
}

// badRangeLoop creates buffer in range loop without sync.Pool
func badRangeLoop() {
	items := []string{"a", "b", "c"}
	// Loop processes items
	for _, item := range items {
		// Buffer created in each iteration
		buf := make([]byte, 512)
		_ = item
		_ = buf
	}
}

// badInfiniteLoop creates buffer in infinite loop
func badInfiniteLoop() {
	// Infinite loop processing
	for {
		// Buffer allocated every iteration
		b := make([]byte, 128)
		_ = b
		break
	}
}

// badWhileStyle creates buffer in while-style loop
func badWhileStyle() {
	counter := 0
	// While-style loop
	for counter < 100 {
		// Buffer created each iteration
		data := make([]byte, 256)
		_ = data
		counter++
	}
}

// BUFFER_SIZE defines buffer size
const BUFFER_SIZE int = 2048

// badConstSizeBuffer creates buffer with const size in loop
func badConstSizeBuffer() {
	// Loop processes items
	for i := 0; i < 50; i++ {
		// Buffer with const size
		buf := make([]byte, BUFFER_SIZE)
		_ = buf
	}
}

// badNestedLoop creates buffer in nested loop
func badNestedLoop() {
	// Outer loop
	for i := 0; i < 10; i++ {
		// Inner loop
		for j := 0; j < 10; j++ {
			// Buffer allocated in nested loop
			temp := make([]byte, 64)
			_ = temp
		}
	}
}

// badMakeWithCapacity creates buffer with capacity in loop
func badMakeWithCapacity() {
	// Loop processes items
	for i := 0; i < 100; i++ {
		// Buffer with both length and capacity
		buffer := make([]byte, 512, 1024)
		_ = buffer
	}
}
