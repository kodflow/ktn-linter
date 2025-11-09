package test007

import "testing"

// TestGoodFunction teste sans utiliser Skip
func TestGoodFunction(t *testing.T) {
	result := GoodFunction()
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

// TestAnotherGoodTest teste sans Skip
func TestAnotherGoodTest(t *testing.T) {
	// Test normal sans skip
	if 1+1 != 2 {
		t.Fatal("Math broken")
	}
}
