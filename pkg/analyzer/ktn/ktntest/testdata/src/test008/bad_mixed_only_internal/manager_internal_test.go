package test008

import "testing"

// Testinitialize teste initialize (fonction privée)
func Testinitialize(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			result := initialize()
			if result != "initialized" {
				t.Errorf("Expected 'initialized', got %s", result)
			}

		})
	}
}
