package ktnstruct_test

import (
	"testing"

	ktnstruct "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestStruct004 vérifie la limite de structs par fichier.
//
// Params:
//   - t: contexte de test
func TestStruct004(t *testing.T) {
	// good.go: 0 errors (1 struct), bad.go: 2 errors (3 structs - les 2 dernières sont en violation)
	tests := []struct {
		name     string
		analyzer string
		expected int
	}{
		{
			name:     "struct004_good_single_struct",
			analyzer: "struct004",
			expected: 2,
		},
		{
			name:     "struct004_verify_analyzer",
			analyzer: "struct004",
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Vérifie que bad.go génère exactement 2 erreurs
			testhelper.TestGoodBad(t, ktnstruct.Analyzer004, tt.analyzer, tt.expected)
		})
	}
}
