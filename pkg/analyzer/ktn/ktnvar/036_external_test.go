package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar036 tests the detection of manual index search patterns.
// Expected errors in bad.go:
// - badIndexOf: manual index search pattern (1)
// - badFindIndex: manual index search pattern (1)
// Total: 2 errors
func TestVar036(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Manual index search patterns that should use slices.Index",
			analyzer:       ktnvar.Analyzer036,
			testdataDir:    "var036",
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
