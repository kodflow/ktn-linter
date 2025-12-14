package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar014 vérifie que les variables de package sont déclarées après les constantes.
func TestVar014(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Vars declared before consts",
			analyzer:       ktnvar.Analyzer014,
			testdataDir:    "var014",
			expectedErrors: 1,
		},
		{
			name:           "Valid vars after consts",
			analyzer:       ktnvar.Analyzer014,
			testdataDir:    "var014",
			expectedErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
