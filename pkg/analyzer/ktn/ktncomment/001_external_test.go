// External tests for ktncomment Analyzer001.
package ktncomment_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
	"golang.org/x/tools/go/analysis"
)

// TestComment001 tests the Analyzer001 for inline comments exceeding 80 characters.
//
// Params:
//   - t: testing context
func TestComment001(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataPath   string
		expectedErrors int
	}{
		{
			name:           "inline comments exceeding 80 characters",
			analyzer:       ktncomment.Analyzer001,
			testdataPath:   "comment001",
			expectedErrors: 4,
		},
	}

	// Parcourir les tests
	for _, tt := range tests {
		// Exécuter chaque test
		t.Run(tt.name, func(t *testing.T) {
			// good.go: 0 errors, bad.go: 4 errors for long inline comments
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataPath, tt.expectedErrors)
		})
	}
}
