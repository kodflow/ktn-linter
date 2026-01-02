package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar035 tests the detection of manual contains patterns.
// Expected errors in bad.go:
// - badContains: manual contains pattern for strings (1)
// - badHasItem: manual contains pattern for ints (1)
// Total: 2 errors
func TestVar035(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Manual contains patterns that should use slices.Contains",
			analyzer:       ktnvar.Analyzer035,
			testdataDir:    "var035",
			expectedErrors: 2,
		},
		{
			name:           "Valid slices.Contains usage and non-contains patterns",
			analyzer:       ktnvar.Analyzer035,
			testdataDir:    "var035",
			expectedErrors: 2,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
