package test005_test

import "testing"

// TestWithTableDriven utilise table-driven tests (BIEN)
func TestWithTableDriven(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"case 1", true},
		{"case 2", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test logic
		})
	}
}
