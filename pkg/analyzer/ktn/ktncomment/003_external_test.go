// External tests for ktncomment Analyzer003.
package ktncomment_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
	"golang.org/x/tools/go/analysis"
)

// TestComment003 tests the Analyzer003 for constant comment requirement.
func TestComment003(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataPath   string
		expectedErrors int
	}{
		{
			name:           "constant comment requirement",
			analyzer:       ktncomment.Analyzer003,
			testdataPath:   "comment003",
			expectedErrors: 22,
		},
	}

	// Parcourir les tests
	for _, tt := range tests {
		// Exécuter chaque test
		t.Run(tt.name, func(t *testing.T) {
			// good.go: 0 errors, bad.go: 22 errors for missing constant comments
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataPath, tt.expectedErrors)
		})
	}
}
