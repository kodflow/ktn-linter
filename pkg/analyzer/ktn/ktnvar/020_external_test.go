package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar020 verifie la detection des slices vides non-nil.
// Erreurs attendues dans bad.go:
// - items := []string{}           (1)
// - data := make([]int, 0)        (1)
// - list := []int{}               (1)
// Total: 3 erreurs
func TestVar020(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Empty slice literals",
			analyzer:       ktnvar.Analyzer020,
			testdataDir:    "var020",
			expectedErrors: 3,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
