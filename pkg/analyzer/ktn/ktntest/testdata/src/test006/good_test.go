package test006

import "testing"

// TestCalculateSum tests the CalculateSum function
func TestCalculateSum(t *testing.T) {
	// Test sum calculation
	result := CalculateSum(2, 3)
	// Verify result
	if result != 5 {
		t.Errorf("Expected 5, got %d", result)
	}
}
