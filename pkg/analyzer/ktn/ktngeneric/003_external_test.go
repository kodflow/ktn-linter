package ktngeneric_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktngeneric "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktngeneric"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestGeneric003(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "deprecated x/exp/constraints import",
			analyzer:       ktngeneric.Analyzer003,
			testdataDir:    "generic003",
			expectedErrors: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// good.go: 0 errors (uses cmp package or no constraints)
			// bad.go: 1 error (uses deprecated x/exp/constraints)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
