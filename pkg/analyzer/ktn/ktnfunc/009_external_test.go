package ktnfunc_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestFunc009 teste KTN-FUNC-009.
//
// Params:
//   - t: contexte de test
func TestFunc009(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataFolder string
		expectedErrors int
	}{
		{
			name:           "func009 with 9 errors",
			analyzer:       ktnfunc.Analyzer009,
			testdataFolder: "func009",
			expectedErrors: 9,
		},
		{
			name:           "func009 consistency check",
			analyzer:       ktnfunc.Analyzer009,
			testdataFolder: "func009",
			expectedErrors: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataFolder, tt.expectedErrors)
		})
	}
}
