package test011

import "testing"

// TestPrivateHelper teste la fonction privée (correct: package test011)
func TestPrivateHelper(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			result := privateHelper()
			if result != "helper" {
				t.Errorf("expected 'helper', got '%s'", result)
			}

		})
	}
}
