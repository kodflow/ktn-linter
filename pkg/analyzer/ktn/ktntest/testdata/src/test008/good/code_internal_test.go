package test008

import "testing"

// TestPrivateFunction teste une fonction privée (white-box testing)
func TestPrivateFunction(t *testing.T) {
	result := privateFunction()
	// Vérification de la condition
	if result != "private" {
		t.Errorf("expected 'private', got '%s'", result)
	}
}
