package ktnstruct_test

import (
	"testing"

	ktnstruct "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestStruct002 vérifie la détection des constructeurs manquants.
func TestStruct002(t *testing.T) {
	// good.go: 0 errors (constructeur NewX présent), bad.go: 1 error (constructeur manquant)
	tests := []struct {
		name     string
		analyzer string
		expected int
	}{
		{
			name:     "struct002_good_constructor_present",
			analyzer: "struct002",
			expected: 1,
		},
		{
			name:     "struct002_verify_analyzer",
			analyzer: "struct002",
			expected: 1,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Vérifie que bad.go génère exactement 1 erreur
			testhelper.TestGoodBad(t, ktnstruct.Analyzer002, tt.analyzer, tt.expected)
		})
	}
}
