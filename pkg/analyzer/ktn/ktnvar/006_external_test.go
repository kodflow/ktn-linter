package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar006(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Builder/Buffer without Grow",
			analyzer:       ktnvar.Analyzer006,
			testdataDir:    "var006",
			expectedErrors: 4,
		},
		{
			name:           "Valid Builder/Buffer with Grow",
			analyzer:       ktnvar.Analyzer006,
			testdataDir:    "var006",
			expectedErrors: 4,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// 4 Builder/Buffer declarations without Grow
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
