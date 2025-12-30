package ktnfunc_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestFunc004 teste KTN-FUNC-004.
func TestFunc004(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataFolder string
		expectedErrors int
	}{
		{
			name:           "func004 with 3 unused private functions",
			analyzer:       ktnfunc.Analyzer004,
			testdataFolder: "func004",
			expectedErrors: 3,
		},
		{
			name:           "func004_special with main, init and callbacks",
			analyzer:       ktnfunc.Analyzer004,
			testdataFolder: "func004_special",
			expectedErrors: 1,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Expected errors depend on test case:
			// - func004: validateTagName, unusedHelper, formatData (3 erreurs)
			// - func004_special: deadFunction (1 erreur)
			// good.go should have no errors:
			// - main() is ignored (entry point)
			// - init() is ignored (called automatically)
			// - run() is used as callback (RunE: run)
			// - helper() is assigned to variable (var helperFunc = helper)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataFolder, tt.expectedErrors)
		})
	}
}
