// External tests for ktncomment Analyzer004.
package ktncomment_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
	"golang.org/x/tools/go/analysis"
)

// TestComment004 tests the Analyzer004 for variable comment requirement.
func TestComment004(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataPath   string
		expectedErrors int
	}{
		{
			name:           "variable comment requirement",
			analyzer:       ktncomment.Analyzer004,
			testdataPath:   "comment004",
			expectedErrors: 20,
		},
	}

	// Parcourir les tests
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Exécuter chaque test
		t.Run(tt.name, func(t *testing.T) {
			// good.go: 0 errors, bad.go: 20 errors for missing variable comments
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataPath, tt.expectedErrors)
		})
	}
}
