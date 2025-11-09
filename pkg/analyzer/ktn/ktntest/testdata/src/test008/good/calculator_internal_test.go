package test008

import "testing"

// TestPrivateHelper teste la fonction privée (white-box)
func TestPrivateHelper(t *testing.T) {
	result := privateHelper()
	if result != 42 {
		t.Errorf("expected 42, got %d", result)
	}
}
