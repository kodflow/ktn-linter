package test008_test

import (
	"testing"

	"test008"
)

// TestPublicFunction teste l'API publique (black-box testing)
func TestPublicFunction(t *testing.T) {
	result := test008.PublicFunction()
	// Vérification de la condition
	if result != "public" {
		t.Errorf("expected 'public', got '%s'", result)
	}
}
