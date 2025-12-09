package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar005(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Make calls with length > 0",
			analyzer:       ktnvar.Analyzer005,
			testdataDir:    "var005",
			expectedErrors: 8,
		},
		{
			name:           "Valid make calls with zero length",
			analyzer:       ktnvar.Analyzer005,
			testdataDir:    "var005",
			expectedErrors: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 8 make calls with length > 0 (VAR-016 cases excluded)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
