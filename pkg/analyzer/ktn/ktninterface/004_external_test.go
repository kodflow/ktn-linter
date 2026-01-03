// Package ktninterface_test provides tests for KTN-INTERFACE-004.
package ktninterface_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktninterface"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
	"golang.org/x/tools/go/analysis"
)

// TestAnalyzer004 tests the empty interface overuse analyzer.
func TestAnalyzer004(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testDir        string
		expectedErrors int
	}{
		{
			name:           "KTN-INTERFACE-004",
			analyzer:       ktninterface.Analyzer004,
			testDir:        "interface004",
			expectedErrors: 5,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			testhelper.TestGoodBad(t, tt.analyzer, tt.testDir, tt.expectedErrors)
		})
	}
}
