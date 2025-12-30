// External tests for ktnreturn Analyzer001.
package ktnreturn_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnreturn"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestReturn001 tests the Analyzer001 for nil returns in slice/map types.
func TestReturn001(t *testing.T) {
	tests := []struct {
		name             string
		analyzer         *analysis.Analyzer
		testDataDir      string
		expectedBadCount int
	}{
		{
			name:             "nil slice returns detection",
			analyzer:         ktnreturn.Analyzer001,
			testDataDir:      "return001",
			expectedBadCount: 6,
		},
		{
			name:             "nil map returns validation",
			analyzer:         ktnreturn.Analyzer001,
			testDataDir:      "return001",
			expectedBadCount: 6,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// good.go: 0 errors, bad.go: 6 errors for nil slice/map returns
			testhelper.TestGoodBad(t, tt.analyzer, tt.testDataDir, tt.expectedBadCount)
		})
	}
}
