package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar018(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Variables with snake_case naming",
			analyzer:       ktnvar.Analyzer018,
			testdataDir:    "var018",
			expectedErrors: 8,
		},
		{
			name:           "Valid camelCase variable naming",
			analyzer:       ktnvar.Analyzer018,
			testdataDir:    "var018",
			expectedErrors: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 8 variables using snake_case (with underscores)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
