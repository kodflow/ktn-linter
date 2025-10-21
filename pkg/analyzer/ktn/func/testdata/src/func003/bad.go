package func003

// badProcessSixItems uses magic number 6
func badProcessSixItems() {
	items := make([]int, 6)
	_ = items
}

// badValidateAge uses magic numbers 18 and 120
func badValidateAge(age int) bool {
	return age >= 18 && age <= 120
}

// badCalculateDiscount uses magic number 0.15
func badCalculateDiscount(price float64) float64 {
	return price * 0.15
}

// badProcessRetries uses magic number 3
func badProcessRetries() {
	maxRetries := 3
	_ = maxRetries
}

// badSetBufferSize uses magic number 1024
func badSetBufferSize() {
	buffer := make([]byte, 1024)
	_ = buffer
}

// badCheckThreshold uses magic number 100
func badCheckThreshold(value int) bool {
	return value > 100
}

// badWaitTimeout uses magic number 30
func badWaitTimeout() {
	timeout := 30
	_ = timeout
}

// badSetPort uses magic number 8080
func badSetPort() {
	port := 8080
	_ = port
}
