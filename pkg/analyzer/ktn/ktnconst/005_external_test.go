package ktnconst_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnconst "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestConst005(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "constants with names exceeding 30 characters",
			analyzer:       ktnconst.Analyzer005,
			testdataDir:    "const005",
			expectedErrors: 5,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// good.go: 0 errors (all names <= 30 chars)
			// bad.go: 5 errors (names exceeding 30 characters)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
