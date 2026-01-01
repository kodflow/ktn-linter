package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar034 verifie la detection du pattern wg.Add(1) + go func() + defer wg.Done().
// Erreurs attendues dans bad.go:
// - badConcurrent: 1 erreur (go func dans boucle)
// - badSimple: 1 erreur (go func simple)
// - badMultiple: 2 erreurs (2x go func)
// Total: 4 erreurs
func TestVar034(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "WaitGroup.Go pattern detection",
			analyzer:       ktnvar.Analyzer034,
			testdataDir:    "var034",
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
