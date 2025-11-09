package test008_test

import (
	"test008"
	"testing"
)

// TestAdd teste la fonction publique Add (black-box)
func TestAdd(t *testing.T) {
	result := test008.Add(2, 3)
	if result != 5 {
		t.Errorf("expected 5, got %d", result)
	}
}
