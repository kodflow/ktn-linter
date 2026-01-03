// Package ktninterface_test provides tests for KTN-INTERFACE-003.
package ktninterface_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktninterface"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
	"golang.org/x/tools/go/analysis"
)

// TestAnalyzer003 tests the single-method interface naming convention analyzer.
func TestAnalyzer003(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testDir        string
		expectedErrors int
	}{
		{
			name:           "KTN-INTERFACE-003",
			analyzer:       ktninterface.Analyzer003,
			testDir:        "interface003",
			expectedErrors: 3,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			testhelper.TestGoodBad(t, tt.analyzer, tt.testDir, tt.expectedErrors)
		})
	}
}
