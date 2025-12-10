// External tests for ktnreturn Analyzer002.
package ktnreturn_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnreturn"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestReturn002 tests the Analyzer002 for nil returns in slice/map types.
func TestReturn002(t *testing.T) {
	tests := []struct {
		name          string
		analyzer      *analysis.Analyzer
		testDataDir   string
		expectedBadCount int
	}{
		{
			name:          "nil slice returns detection",
			analyzer:      ktnreturn.Analyzer002,
			testDataDir:   "return002",
			expectedBadCount: 6,
		},
		{
			name:          "nil map returns validation",
			analyzer:      ktnreturn.Analyzer002,
			testDataDir:   "return002",
			expectedBadCount: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// good.go: 0 errors, bad.go: 6 errors for nil slice/map returns
			testhelper.TestGoodBad(t, tt.analyzer, tt.testDataDir, tt.expectedBadCount)
		})
	}
}
