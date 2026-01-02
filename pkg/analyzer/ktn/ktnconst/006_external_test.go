package ktnconst_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnconst "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestConst006(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "constants shadowing built-in identifiers",
			analyzer:       ktnconst.Analyzer006,
			testdataDir:    "const006",
			expectedErrors: 15,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// good.go: 0 errors (no shadowing of built-ins)
			// bad.go: 15 errors (constants shadowing built-in identifiers)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
