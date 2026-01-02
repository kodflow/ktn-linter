package ktngeneric_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktngeneric "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktngeneric"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestGeneric001(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "generic functions without comparable constraint",
			analyzer:       ktngeneric.Analyzer001,
			testdataDir:    "generic001",
			expectedErrors: 2,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// good.go: 0 errors (proper comparable constraint)
			// bad.go: 2 errors (badContains and badIndex using == without comparable)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
