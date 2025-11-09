package test009_test

import (
	"test009"
	"testing"
)

// TestPublicFunction teste la fonction publique (correct dans external)
func TestPublicFunction(t *testing.T) {
	result := test009.PublicFunction()
	if result != "public" {
		t.Errorf("expected 'public', got '%s'", result)
	}
}
