package ktnstruct_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestStruct006 vérifie la détection des champs privés tagués dans les DTOs.
//
// Params:
//   - t: contexte de test
func TestStruct006(t *testing.T) {
	// 4 champs privés avec tags dans des DTOs
	tests := []struct {
		name     string
		analyzer string
		expected int
	}{
		{
			name:     "struct006_good_no_private_tagged_fields",
			analyzer: "struct006",
			expected: 4,
		},
		{
			name:     "struct006_verify_analyzer",
			analyzer: "struct006",
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Vérifie que bad.go génère exactement 4 erreurs
			testhelper.TestGoodBad(t, ktnstruct.Analyzer006, tt.analyzer, tt.expected)
		})
	}
}
