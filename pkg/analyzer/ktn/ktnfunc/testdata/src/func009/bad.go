// Package func009 contains test cases for KTN rules.
package func009

// badProcessSixItems creates a slice with magic number 6 (violates KTN-FUNC-009).
func badProcessSixItems() {
	// Create slice with magic number 6 (should be a constant)
	items := [6]int{}
	_ = items
}

// badValidateAge checks age validity using magic numbers 18 and 120 (violates KTN-FUNC-009).
//
// Params:
//   - age: the age to validate
//
// Returns:
//   - bool: true if age is between 18 and 120, false otherwise
func badValidateAge(age int) bool {
	// Check age range with magic numbers 18 and 120 (should be constants)
	if age >= 18 && age <= 120 {
		// Age is in valid range
		return true
	}
	// Age is out of valid range
	return false
}

// badCalculateDiscount applies discount using magic number 0.15 (violates KTN-FUNC-009).
//
// Params:
//   - price: the original price
//
// Returns:
//   - float64: the discount amount
func badCalculateDiscount(price float64) float64 {
	// Calculate discount with magic number 0.15 (should be a constant)
	return price * 0.15
}

// badProcessRetries sets max retries using magic number 3 (violates KTN-FUNC-009).
func badProcessRetries() {
	// Set retries with magic number 3 (should be a constant)
	maxRetries := 3
	_ = maxRetries
}

// badSetBufferSize creates a buffer using magic number 1024 (violates KTN-FUNC-009).
func badSetBufferSize() {
	// Create buffer with magic number 1024 (should be a constant)
	buffer := [1024]byte{}
	_ = buffer
}

// badCheckThreshold verifies value against magic number 100 (violates KTN-FUNC-009).
//
// Params:
//   - value: the value to check
//
// Returns:
//   - bool: true if value exceeds 100, false otherwise
func badCheckThreshold(value int) bool {
	// Check threshold with magic number 100 (should be a constant)
	if value > 100 {
		// Value exceeds threshold
		return true
	}
	// Value is below or equal to threshold
	return false
}

// badWaitTimeout sets a timeout using magic number 30 (violates KTN-FUNC-009).
func badWaitTimeout() {
	// Set timeout with magic number 30 (should be a constant)
	timeout := 30
	_ = timeout
}

// badSetPort assigns a port using magic number 8080 (violates KTN-FUNC-009).
func badSetPort() {
	// Set port with magic number 8080 (should be a constant)
	port := 8080
	_ = port
}

// init utilise les fonctions privées
func init() {
	// Appel de badProcessSixItems
	badProcessSixItems()
	// Appel de badValidateAge
	_ = badValidateAge(0)
	// Appel de badCalculateDiscount
	_ = badCalculateDiscount(0)
	// Appel de badProcessRetries
	badProcessRetries()
	// Appel de badSetBufferSize
	badSetBufferSize()
	// Appel de badCheckThreshold
	_ = badCheckThreshold(0)
	// Appel de badWaitTimeout
	badWaitTimeout()
	// Appel de badSetPort
	badSetPort()
}
