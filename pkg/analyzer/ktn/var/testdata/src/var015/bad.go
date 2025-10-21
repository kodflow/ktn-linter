package var015

// Bad: Creating []byte buffers repeatedly in loops without sync.Pool

// badProcessInLoop creates buffer in loop without sync.Pool
func badProcessInLoop() {
	// Loop processes items
	for i := 0; i < 100; i++ {
		// Buffer created repeatedly in loop
		buffer := make([]byte, 1024) // want "KTN-VAR-015"
		_ = buffer
	}
}

// badRangeLoop creates buffer in range loop without sync.Pool
func badRangeLoop() {
	items := []string{"a", "b", "c"}
	// Loop processes items
	for _, item := range items {
		// Buffer created in each iteration
		buf := make([]byte, 512) // want "KTN-VAR-015"
		_ = item
		_ = buf
	}
}

// badInfiniteLoop creates buffer in infinite loop
func badInfiniteLoop() {
	// Infinite loop processing
	for {
		// Buffer allocated every iteration
		b := make([]byte, 128) // want "KTN-VAR-015"
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
		data := make([]byte, 256) // want "KTN-VAR-015"
		_ = data
		counter++
	}
}
