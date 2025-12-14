package ktnstruct_test

import (
	"testing"

	ktnstruct "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestStruct001 vérifie la détection des structs sans interface.
func TestStruct001(t *testing.T) {
	// good.go: 0 errors (interface complète)
	// bad.go: 2 errors:
	//   - BadUserService: struct sans interface
	//   - BadIncompleteImpl: compile-time check présent mais interface incomplète
	tests := []struct {
		name     string
		analyzer string
		expected int
	}{
		{
			name:     "struct001_good_interface_complete",
			analyzer: "struct001",
			expected: 2,
		},
		{
			name:     "struct001_verify_analyzer",
			analyzer: "struct001",
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Vérifie que bad.go génère exactement 1 erreur
			testhelper.TestGoodBad(t, ktnstruct.Analyzer001, tt.analyzer, tt.expected)
		})
	}
}
