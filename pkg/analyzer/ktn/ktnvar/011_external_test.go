package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar011 vérifie la détection du shadowing de variables.
func TestVar011(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Variable shadowing detected",
			analyzer:       ktnvar.Analyzer011,
			testdataDir:    "var011",
			expectedErrors: 5,
		},
		{
			name:           "No variable shadowing",
			analyzer:       ktnvar.Analyzer011,
			testdataDir:    "var011",
			expectedErrors: 5,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// 5 cas de shadowing attendus
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
