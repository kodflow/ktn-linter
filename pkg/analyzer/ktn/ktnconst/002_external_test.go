package ktnconst_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnconst "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestConst002(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "scattered constants",
			analyzer:       ktnconst.Analyzer002,
			testdataDir:    "const002",
			expectedErrors: 7,
		},
		{
			name:           "valid grouped constants",
			analyzer:       ktnconst.Analyzer002,
			testdataDir:    "const002",
			expectedErrors: 7,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// good.go: 0 errors
			// bad.go: 7 errors (4 scattered + 2 after type + 1 after func)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
