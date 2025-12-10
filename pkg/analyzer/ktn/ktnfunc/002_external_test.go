package ktnfunc_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestFunc002 teste KTN-FUNC-002.
//
// Params:
//   - t: contexte de test
func TestFunc002(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataFolder string
		expectedErrors int
	}{
		{
			name:           "func002 with 4 errors",
			analyzer:       ktnfunc.Analyzer002,
			testdataFolder: "func002",
			expectedErrors: 4,
		},
		{
			name:           "func002 consistency check",
			analyzer:       ktnfunc.Analyzer002,
			testdataFolder: "func002",
			expectedErrors: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataFolder, tt.expectedErrors)
		})
	}
}
