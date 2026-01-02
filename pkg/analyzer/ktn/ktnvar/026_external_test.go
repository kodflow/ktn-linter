package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar026 verifie la detection des patterns min/max manuels.
// Erreurs attendues dans bad.go:
// - badMaxFloat: math.Max (1)
// - badMinFloat: math.Min (1)
// Total: 2 erreurs (math.Min et math.Max seulement)
func TestVar026(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "min/max built-in detection",
			analyzer:       ktnvar.Analyzer026,
			testdataDir:    "var026",
			expectedErrors: 2,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
