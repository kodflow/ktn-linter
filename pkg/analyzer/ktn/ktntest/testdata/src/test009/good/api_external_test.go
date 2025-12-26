package test011_test

import (
	"test011"
	"testing"
)

// TestPublicFunc teste la fonction publique (correct: package test011_test)
func TestPublicFunc(t *testing.T) {
	result := test011.PublicFunc()
	if result != "helper" {
		t.Errorf("expected 'helper', got '%s'", result)
	}
}
