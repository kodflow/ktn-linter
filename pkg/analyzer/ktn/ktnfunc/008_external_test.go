package ktnfunc_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestFunc008 teste KTN-FUNC-008.
func TestFunc008(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataFolder string
		expectedErrors int
	}{
		{
			name:           "func008 with 10 unused parameters",
			analyzer:       ktnfunc.Analyzer008,
			testdataFolder: "func008",
			expectedErrors: 10,
		},
		{
			name:           "func008 consistency check",
			analyzer:       ktnfunc.Analyzer008,
			testdataFolder: "func008",
			expectedErrors: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Expected errors in bad.go:
			// - Delete: ctx (_ = ctx), req, resp (3 erreurs)
			// - ProcessData: ctx, options (2 erreurs)
			// - PartialIgnore: a, b (_ = b), c (3 erreurs)
			// - BadHandler.Handle: ctx, data (2 erreurs - méthode implémentant interface)
			// Total: 10 erreurs
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataFolder, tt.expectedErrors)
		})
	}
}
