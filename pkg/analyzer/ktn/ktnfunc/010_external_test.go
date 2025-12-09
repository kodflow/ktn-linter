package ktnfunc_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestFunc010 teste KTN-FUNC-010.
//
// Params:
//   - t: contexte de test
func TestFunc010(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataFolder string
		expectedErrors int
	}{
		{
			name:           "func010 with 4 naked return errors",
			analyzer:       ktnfunc.Analyzer010,
			testdataFolder: "func010",
			expectedErrors: 4,
		},
		{
			name:           "func010 consistency check",
			analyzer:       ktnfunc.Analyzer010,
			testdataFolder: "func010",
			expectedErrors: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// func010/bad.go doit avoir 4 erreurs (naked returns dans fonctions trop longues)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataFolder, tt.expectedErrors)
		})
	}
}
