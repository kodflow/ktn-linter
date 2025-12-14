package ktnstruct_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestStruct007 vérifie la convention de nommage des getters/setters.
// Les getters sont OPTIONNELS, mais s'ils existent, ils doivent suivre la convention.
// Note: La détection du préfixe Get est gérée par STRUCT-003.
// STRUCT-007 vérifie uniquement le mismatch nom getter vs champ retourné.
// Erreurs attendues dans bad.go:
// - Value() retourne le champ 'data', devrait être nommé Data() (1)
// Total: 1 erreur
func TestStruct007(t *testing.T) {
	// 1 violation: mismatch getter/champ
	tests := []struct {
		name     string
		analyzer string
		expected int
	}{
		{
			name:     "struct007_good_getter_naming",
			analyzer: "struct007",
			expected: 1,
		},
		{
			name:     "struct007_verify_analyzer",
			analyzer: "struct007",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Vérifie que bad.go génère exactement 1 erreur
			testhelper.TestGoodBad(t, ktnstruct.Analyzer007, tt.analyzer, tt.expected)
		})
	}
}
