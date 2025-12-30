package ktnfunc_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestFunc006 teste KTN-FUNC-006.
func TestFunc006(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataFolder string
		expectedErrors int
	}{
		{
			name:           "func006 with 10 errors",
			analyzer:       ktnfunc.Analyzer006,
			testdataFolder: "func006",
			expectedErrors: 10,
		},
		{
			name:           "func006 consistency check",
			analyzer:       ktnfunc.Analyzer006,
			testdataFolder: "func006",
			expectedErrors: 10,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataFolder, tt.expectedErrors)
		})
	}
}
