// External tests for ktninterface Analyzer001.
package ktninterface_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktninterface"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestInterface001 tests the Analyzer001 for unused interface declarations.
func TestInterface001(t *testing.T) {
	tests := []struct {
		name             string
		analyzer         *analysis.Analyzer
		testDataDir      string
		expectedBadCount int
	}{
		{
			name:             "unused interfaces detection",
			analyzer:         ktninterface.Analyzer001,
			testDataDir:      "interface001",
			expectedBadCount: 3,
		},
		{
			name:             "interface usage validation",
			analyzer:         ktninterface.Analyzer001,
			testDataDir:      "interface001",
			expectedBadCount: 3,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// good.go: 0 errors, bad.go: 3 errors for unused interfaces
			testhelper.TestGoodBad(t, tt.analyzer, tt.testDataDir, tt.expectedBadCount)
		})
	}
}
