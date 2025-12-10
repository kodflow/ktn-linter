// External tests for ktncomment Analyzer002.
package ktncomment_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
	"golang.org/x/tools/go/analysis"
)

// TestComment002 tests the Analyzer002 for package comment requirement.
//
// Params:
//   - t: testing context
func TestComment002(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataPath   string
		expectedErrors int
	}{
		{
			name:           "package comment requirement",
			analyzer:       ktncomment.Analyzer002,
			testdataPath:   "comment002",
			expectedErrors: 2,
		},
	}

	// Parcourir les tests
	for _, tt := range tests {
		// Exécuter chaque test
		t.Run(tt.name, func(t *testing.T) {
			// good/: 0 errors, bad/: 2 errors for missing package comments
			testhelper.TestGoodBadPackage(t, tt.analyzer, tt.testdataPath, tt.expectedErrors)
		})
	}
}
