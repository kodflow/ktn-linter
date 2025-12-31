package ktnconst_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnconst "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestConst001(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "constants without explicit type",
			analyzer:       ktnconst.Analyzer001,
			testdataDir:    "const001",
			expectedErrors: 47,
		},
		{
			name:           "valid constants with explicit type",
			analyzer:       ktnconst.Analyzer001,
			testdataDir:    "const001",
			expectedErrors: 47,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// good.go: 0 errors
			// bad.go: 47 errors (all constant types without explicit type)
			// Sections: basic(8) + literals(8) + rune(7) + complex(3) + string(4)
			//         + iota(2) + expr(6) + multi(4) + edge(5) = 47
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
