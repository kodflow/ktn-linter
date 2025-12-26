package ktnconst_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnconst "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestConst003(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "constants with underscores violating CamelCase",
			analyzer:       ktnconst.Analyzer003,
			testdataDir:    "const003",
			expectedErrors: 20,
		},
		{
			name:           "valid CamelCase constants",
			analyzer:       ktnconst.Analyzer003,
			testdataDir:    "const003",
			expectedErrors: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// good.go: 0 errors (CamelCase is valid)
			// bad.go: 20 errors (constants with underscores violate CamelCase convention)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
