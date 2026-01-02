package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar024 verifie la detection de interface{} au lieu de any.
// Erreurs attendues dans bad.go:
// - badProcess: parametre interface{} (1)
// - badX: variable interface{} (1)
// - BadContainer.value: champ interface{} (1)
// - badReturns: retour interface{} (1)
// Total: 4 erreurs
func TestVar024(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "interface{} instead of any",
			analyzer:       ktnvar.Analyzer024,
			testdataDir:    "var024",
			expectedErrors: 4,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
