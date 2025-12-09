package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar010(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Buffers created in loops",
			analyzer:       ktnvar.Analyzer010,
			testdataDir:    "var010",
			expectedErrors: 8,
		},
		{
			name:           "Valid buffer reuse outside loops",
			analyzer:       ktnvar.Analyzer010,
			testdataDir:    "var010",
			expectedErrors: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 8 buffers créés dans des boucles (4 original + 4 nouveaux edge cases)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
