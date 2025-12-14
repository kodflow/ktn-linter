package ktnfunc_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestFunc005 teste KTN-FUNC-005.
func TestFunc005(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataFolder string
		expectedErrors int
	}{
		{
			name:           "func005 with 2 too long functions",
			analyzer:       ktnfunc.Analyzer005,
			testdataFolder: "func005",
			expectedErrors: 2,
		},
		{
			name:           "func005 consistency check",
			analyzer:       ktnfunc.Analyzer005,
			testdataFolder: "func005",
			expectedErrors: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// func005/bad.go doit avoir 2 erreurs (TooLong: 36 lignes, VeryLong: 37 lignes)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataFolder, tt.expectedErrors)
		})
	}
}
