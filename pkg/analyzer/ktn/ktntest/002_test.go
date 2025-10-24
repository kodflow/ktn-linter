package ktntest

import (
	"testing"
)

// TestIsExemptTestFile tests isExemptTestFile from 002.go
func TestIsExemptTestFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{"helper test", "helper_test.go", true},
		{"integration test", "integration_test.go", true},
		{"suite test", "suite_test.go", true},
		{"main test", "main_test.go", true},
		{"regular test", "foo_test.go", false},
		{"rule test", "001_test.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isExemptTestFile(tt.filename)
			if got != tt.want {
				t.Errorf("isExemptTestFile(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}
