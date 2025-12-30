package ktnstruct_test

import (
	"testing"

	ktnstruct "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestStruct003 vérifie la détection des getters avec préfixe Get.
func TestStruct003(t *testing.T) {
	// good.go: 0 errors (getters idiomatiques sans Get), bad.go: 3 errors (getters avec préfixe Get)
	tests := []struct {
		name     string
		analyzer string
		expected int
	}{
		{
			name:     "struct003_good_idiomatic_getters",
			analyzer: "struct003",
			expected: 3,
		},
		{
			name:     "struct003_verify_analyzer",
			analyzer: "struct003",
			expected: 3,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Vérifie que bad.go génère exactement 3 erreurs
			testhelper.TestGoodBad(t, ktnstruct.Analyzer003, tt.analyzer, tt.expected)
		})
	}
}
