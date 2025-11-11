package good

import "testing"

// Test using Test_functionName convention for private functions
func Test_populateData(t *testing.T) {
	result := populateData("test")
	if result != "populated: test" {
		t.Errorf("expected 'populated: test', got '%s'", result)
	}
}

// Test using Test_functionName convention for another private function
func Test_validateInput(t *testing.T) {
	if !validateInput("valid") {
		t.Error("expected true for non-empty input")
	}
	if validateInput("") {
		t.Error("expected false for empty input")
	}
}
