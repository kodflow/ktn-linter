// External tests for ktncomment Analyzer006.
package ktncomment_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
	"golang.org/x/tools/go/analysis"
)

// TestComment006 tests the Analyzer006 for function documentation format.
//
// Params:
//   - t: testing context
func TestComment006(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataPath   string
		expectedErrors int
	}{
		{
			name:           "function documentation format",
			analyzer:       ktncomment.Analyzer006,
			testdataPath:   "comment006",
			expectedErrors: 6,
		},
	}

	// Parcourir les tests
	for _, tt := range tests {
		// Exécuter chaque test
		t.Run(tt.name, func(t *testing.T) {
			// good.go: 0 errors, bad.go: 6 errors for missing/invalid function docs
			// - badNoDoc: no doc, badWrongFormat: wrong format, badMissingParams: no Params
			// - badMissingReturns: no Returns, badEmptyParamsSection: empty Params
			// - badEmptyReturnsSection: empty Returns
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataPath, tt.expectedErrors)
		})
	}
}
