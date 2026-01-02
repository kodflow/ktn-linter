package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar029 tests the detection of manual slice grow patterns.
// Expected errors in bad.go:
// - badGrow: manual grow pattern (if cap-len < n, make+copy) (1)
// - badGrowString: manual grow pattern for string slice (1)
// - badGrowBytes: manual grow pattern for byte slice (1)
// Total: 3 errors
func TestVar029(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Manual slice grow patterns that should use slices.Grow",
			analyzer:       ktnvar.Analyzer029,
			testdataDir:    "var029",
			expectedErrors: 3,
		},
		{
			name:           "Valid slices.Grow usage and non-grow patterns",
			analyzer:       ktnvar.Analyzer029,
			testdataDir:    "var029",
			expectedErrors: 3,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
