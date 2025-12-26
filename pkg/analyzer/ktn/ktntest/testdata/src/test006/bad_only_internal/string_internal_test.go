package test008

import "testing"

// TestConcat teste Concat (devrait être dans external)
func TestConcat(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			result := Concat("hello", "world")
			if result != "helloworld" {
				t.Errorf("expected 'helloworld', got '%s'", result)
			}

		})
	}
}
