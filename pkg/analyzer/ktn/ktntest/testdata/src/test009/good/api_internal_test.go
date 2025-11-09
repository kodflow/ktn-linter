package test009

import "testing"

// TestPrivateHelper teste la fonction privée (correct dans internal)
func TestPrivateHelper(t *testing.T) {
	result := privateHelper()
	if result != "helper" {
		t.Errorf("expected 'helper', got '%s'", result)
	}
}
