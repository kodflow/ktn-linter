package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar037 tests the detection of manual map key/value collection patterns.
// Expected errors in bad.go:
// - badGetKeys: manual keys collection (1)
// - badGetValues: manual values collection (1)
// Total: 2 errors
func TestVar037(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Manual map keys collection that should use slices.Collect(maps.Keys())",
			analyzer:       ktnvar.Analyzer037,
			testdataDir:    "var037",
			expectedErrors: 2,
		},
		{
			name:           "Valid maps.Keys/Values usage and non-collection patterns",
			analyzer:       ktnvar.Analyzer037,
			testdataDir:    "var037",
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
