package test008

import "testing"

// TestGetValue teste GetValue dans _internal (mauvais fichier !)
func TestGetValue(t *testing.T) {
	result := GetValue()
	if result != "value" {
		t.Errorf("Expected 'value', got %s", result)
	}
}
