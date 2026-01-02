package ktnvar

import (
	"testing"
)

// TestMaxVarNameLength005 tests the maximum name length constant.
func TestMaxVarNameLength005(t *testing.T) {
	if maxVarNameLength005 != 30 {
		t.Errorf("maxVarNameLength005 = %d, want 30", maxVarNameLength005)
	}
}

// TestNameLengthLimits005 tests edge cases for name length.
func TestNameLengthLimits005(t *testing.T) {
	tests := []struct {
		name     string
		varName  string
		tooLong  bool
	}{
		{
			name:    "exactly 30 chars",
			varName: "ThisIsExactlyThirtyCharacterss",
			tooLong: false,
		},
		{
			name:    "31 chars",
			varName: "ThisIsExactlyThirtyOneCharacter",
			tooLong: true,
		},
		{
			name:    "short name",
			varName: "ok",
			tooLong: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := len(tt.varName) > maxVarNameLength005
			if result != tt.tooLong {
				t.Errorf("len(%q) > %d = %v, want %v",
					tt.varName, maxVarNameLength005, result, tt.tooLong)
			}
		})
	}
}
