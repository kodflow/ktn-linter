package test008

import "testing"

// TestConcat teste Concat (devrait être dans external)
func TestConcat(t *testing.T) {
	result := Concat("hello", "world")
	if result != "helloworld" {
		t.Errorf("expected 'helloworld', got '%s'", result)
	}
}
