package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar017(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Maps without capacity hints",
			analyzer:       ktnvar.Analyzer017,
			testdataDir:    "var017",
			expectedErrors: 7,
		},
		{
			name:           "Valid maps with capacity",
			analyzer:       ktnvar.Analyzer017,
			testdataDir:    "var017",
			expectedErrors: 7,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// 7 maps without capacity hints
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
