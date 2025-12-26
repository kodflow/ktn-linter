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
			expectedErrors: 16,
		},
		{
			name:           "valid constants with explicit type",
			analyzer:       ktnconst.Analyzer001,
			testdataDir:    "const001",
			expectedErrors: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// good.go: 0 errors, bad.go: 16 errors (constants without explicit type)
			// - 8 basic + 5 numeric + 1 rune + 1 iota + 2 multi-name - 1 inherited = 16
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
