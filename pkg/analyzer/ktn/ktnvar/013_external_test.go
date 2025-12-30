package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar013(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Scattered var declarations",
			analyzer:       ktnvar.Analyzer013,
			testdataDir:    "var013",
			expectedErrors: 4,
		},
		{
			name:           "Valid grouped var declarations",
			analyzer:       ktnvar.Analyzer013,
			testdataDir:    "var013",
			expectedErrors: 4,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// 5 scattered var declarations (2 single + 3 groups after first)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
