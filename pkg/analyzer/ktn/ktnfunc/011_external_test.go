package ktnfunc_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestFunc011 teste KTN-FUNC-011.
func TestFunc011(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataFolder string
		expectedErrors int
	}{
		{
			name:           "func011 with 6 errors",
			analyzer:       ktnfunc.Analyzer011,
			testdataFolder: "func011",
			expectedErrors: 6,
		},
		{
			name:           "func011 consistency check",
			analyzer:       ktnfunc.Analyzer011,
			testdataFolder: "func011",
			expectedErrors: 6,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataFolder, tt.expectedErrors)
		})
	}
}
