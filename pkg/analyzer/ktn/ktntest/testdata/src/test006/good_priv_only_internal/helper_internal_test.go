package test008

import "testing"

// TestparseInput teste parseInput
func TestparseInput(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			result := parseInput("test")
			if result != 4 {
				t.Errorf("Expected 4, got %d", result)
			}

		})
	}
}

// TestformatOutput teste formatOutput
func TestformatOutput(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			result := formatOutput("test")
			if result != "[test]" {
				t.Errorf("Expected [test], got %s", result)
			}

		})
	}
}
