package ktnfunc_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestFunc001 teste KTN-FUNC-001.
func TestFunc001(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataFolder string
		expectedErrors int
	}{
		{
			name:           "func001 with 7 errors",
			analyzer:       ktnfunc.Analyzer001,
			testdataFolder: "func001",
			expectedErrors: 7,
		},
		{
			name:           "func001 consistency check",
			analyzer:       ktnfunc.Analyzer001,
			testdataFolder: "func001",
			expectedErrors: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 6 erreurs de position + 1 erreur "2 types error"
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataFolder, tt.expectedErrors)
		})
	}
}
