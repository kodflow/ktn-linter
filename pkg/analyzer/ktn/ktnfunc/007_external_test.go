package ktnfunc_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestFunc007 teste KTN-FUNC-007.
func TestFunc007(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataFolder string
		expectedErrors int
	}{
		{
			name:           "func007 with 6 errors",
			analyzer:       ktnfunc.Analyzer007,
			testdataFolder: "func007",
			expectedErrors: 6,
		},
		{
			name:           "func007 consistency check",
			analyzer:       ktnfunc.Analyzer007,
			testdataFolder: "func007",
			expectedErrors: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataFolder, tt.expectedErrors)
		})
	}
}
