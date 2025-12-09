package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar004(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Slice initialization without capacity",
			analyzer:       ktnvar.Analyzer004,
			testdataDir:    "var004",
			expectedErrors: 8,
		},
		{
			name:           "Valid slice initialization with capacity",
			analyzer:       ktnvar.Analyzer004,
			testdataDir:    "var004",
			expectedErrors: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 8 errors:
			// - 7 make([]T, 0) calls without capacity
			// - 1 []T{} literal followed by append (should use make)
			// []T{} in returns/structs are NOT reported (avoid false positives)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
