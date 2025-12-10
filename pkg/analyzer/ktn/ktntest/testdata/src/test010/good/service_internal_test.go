package test010

import "testing"

// Test_privateImplementation teste la fonction privée (correct dans internal)
func Test_privateImplementation(t *testing.T) {
	result := privateImplementation()
	if result != "impl" {
		t.Errorf("expected 'impl', got '%s'", result)
	}
}
