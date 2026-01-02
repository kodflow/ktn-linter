package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar031(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Manual map cloning pattern",
			analyzer:       ktnvar.Analyzer031,
			testdataDir:    "var031",
			expectedErrors: 3,
		},
		{
			name:           "Valid maps.Clone usage",
			analyzer:       ktnvar.Analyzer031,
			testdataDir:    "var031",
			expectedErrors: 3,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// 3 manual map clone patterns detected
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
