// External tests for ktncomment Analyzer007.
package ktncomment_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
	"golang.org/x/tools/go/analysis"
)

// TestComment007 tests the Analyzer007 for control block comments.
//
// Params:
//   - t: testing context
func TestComment007(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataPath   string
		expectedErrors int
	}{
		{
			name:           "control block comments",
			analyzer:       ktncomment.Analyzer007,
			testdataPath:   "comment007",
			expectedErrors: 37,
		},
	}

	// Parcourir les tests
	for _, tt := range tests {
		// Exécuter chaque test
		t.Run(tt.name, func(t *testing.T) {
			// good.go: 0 errors, bad.go: 37 errors for missing block comments
			// (règle stricte: tous les blocs, returns, et if doivent avoir un commentaire)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataPath, tt.expectedErrors)
		})
	}
}
