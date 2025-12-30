package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar002(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Variables without explicit type",
			analyzer:       ktnvar.Analyzer002,
			testdataDir:    "var002",
			expectedErrors: 8,
		},
		{
			name:           "Valid explicit type declarations",
			analyzer:       ktnvar.Analyzer002,
			testdataDir:    "var002",
			expectedErrors: 8,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// 8 variables without explicit type (zero-values are now valid)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
