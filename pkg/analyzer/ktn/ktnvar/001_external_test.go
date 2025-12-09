package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar001(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "SCREAMING_SNAKE_CASE naming violations",
			analyzer:       ktnvar.Analyzer001,
			testdataDir:    "var001",
			expectedErrors: 9,
		},
		{
			name:           "Valid SCREAMING_SNAKE_CASE naming",
			analyzer:       ktnvar.Analyzer001,
			testdataDir:    "var001",
			expectedErrors: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 9 variables with SCREAMING_SNAKE_CASE naming (6 original + 3 acronym cases)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
