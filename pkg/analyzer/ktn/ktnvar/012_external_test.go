package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar012(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Allocations in loops",
			analyzer:       ktnvar.Analyzer012,
			testdataDir:    "var012",
			expectedErrors: 7,
		},
		{
			name:           "Valid pre-loop allocations",
			analyzer:       ktnvar.Analyzer012,
			testdataDir:    "var012",
			expectedErrors: 7,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// 7 allocations dans des boucles (5 assignements + 2 déclarations var)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
