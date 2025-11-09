package test010

import "testing"

// TestPrivateImplementation teste la fonction privée (correct dans internal)
func TestPrivateImplementation(t *testing.T) {
	result := privateImplementation()
	if result != "impl" {
		t.Errorf("expected 'impl', got '%s'", result)
	}
}
