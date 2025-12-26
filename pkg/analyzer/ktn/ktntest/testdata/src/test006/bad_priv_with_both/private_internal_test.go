package test008

import "testing"

// TestcomputeValue teste computeValue dans _internal (bon fichier !)
func TestcomputeValue(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			result := computeValue(5)
			if result != 10 {
				t.Errorf("Expected 10, got %d", result)
			}

		})
	}
}
