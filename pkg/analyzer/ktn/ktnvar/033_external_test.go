package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar033 tests the detection of patterns that can use cmp.Or.
// Expected errors in bad.go:
// - badGetPort: if port != 0 pattern (1)
// - badGetHost: if host != "" pattern (1)
// - badGetPointer: if ptr != nil pattern (1)
// - badGetSlice: if slice != nil pattern (1)
// Total: 4 errors
func TestVar033(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Patterns that should use cmp.Or",
			analyzer:       ktnvar.Analyzer033,
			testdataDir:    "var033",
			expectedErrors: 4,
		},
		{
			name:           "Valid cmp.Or usage and non-matching patterns",
			analyzer:       ktnvar.Analyzer033,
			testdataDir:    "var033",
			expectedErrors: 4,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
