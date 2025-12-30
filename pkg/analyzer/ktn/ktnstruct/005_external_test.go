package ktnstruct_test

import (
	"testing"

	ktnstruct "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestStruct005 vérifie l'ordre des champs exportés/privés.
func TestStruct005(t *testing.T) {
	// good.go: 0 errors (champs exportés avant privés), bad.go: 5 errors (champs mélangés)
	tests := []struct {
		name     string
		analyzer string
		expected int
	}{
		{
			name:     "struct005_good_exported_before_private",
			analyzer: "struct005",
			expected: 5,
		},
		{
			name:     "struct005_verify_analyzer",
			analyzer: "struct005",
			expected: 5,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Vérifie que bad.go génère exactement 5 erreurs
			testhelper.TestGoodBad(t, ktnstruct.Analyzer005, tt.analyzer, tt.expected)
		})
	}
}
