package ktnfunc_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestFunc003 teste KTN-FUNC-003.
// Erreurs attendues dans bad.go:
// - badCheckPositive: else après return
// - badProcessValue: else après return
// - badFindMax: else après return
// - badLoopExample: else après continue
// - badSwitchExample: else après break
// - badValidateInput: else après return
// - badPanicExample: else après panic
// - badElseIfExample: else if après return + else après return (2 erreurs)
// Total: 9 erreurs
//
// Params:
//   - t: contexte de test
func TestFunc003(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataFolder string
		expectedErrors int
	}{
		{
			name:           "func003 with 9 errors",
			analyzer:       ktnfunc.Analyzer003,
			testdataFolder: "func003",
			expectedErrors: 9,
		},
		{
			name:           "func003 consistency check",
			analyzer:       ktnfunc.Analyzer003,
			testdataFolder: "func003",
			expectedErrors: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataFolder, tt.expectedErrors)
		})
	}
}
