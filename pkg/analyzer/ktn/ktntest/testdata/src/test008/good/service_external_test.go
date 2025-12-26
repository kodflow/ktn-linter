package test010_test

import (
	"test010"
	"testing"
)

// TestPublicService teste la fonction publique (correct dans external)
func TestPublicService(t *testing.T) {
	result := test010.PublicService()
	if result != "impl" {
		t.Errorf("expected 'impl', got '%s'", result)
	}
}
