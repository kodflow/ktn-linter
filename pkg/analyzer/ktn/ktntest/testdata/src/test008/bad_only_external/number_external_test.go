package test008_test

import (
	"test008"
	"testing"
)

// TestDouble teste la fonction publique
func TestDouble(t *testing.T) {
	result := test008.Double(5)
	if result != 10 {
		t.Errorf("expected 10, got %d", result)
	}
}
