package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar003(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Variables using var with initialization",
			analyzer:       ktnvar.Analyzer003,
			testdataDir:    "var003",
			expectedErrors: 15,
		},
		{
			name:           "Valid short declarations with :=",
			analyzer:       ktnvar.Analyzer003,
			testdataDir:    "var003",
			expectedErrors: 15,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// 15 variables using var with initialization instead of := (13 + 2 dans select)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
