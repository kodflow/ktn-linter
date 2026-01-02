package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar011(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "String concatenation in loops",
			analyzer:       ktnvar.Analyzer011,
			testdataDir:    "var011",
			expectedErrors: 6,
		},
		{
			name:           "Valid string building with Builder",
			analyzer:       ktnvar.Analyzer011,
			testdataDir:    "var011",
			expectedErrors: 6,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// 6 string concatenation errors detected
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
