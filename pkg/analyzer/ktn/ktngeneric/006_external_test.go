package ktngeneric_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktngeneric "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktngeneric"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestGeneric006(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "generic functions using ordered/arithmetic ops without cmp.Ordered",
			analyzer:       ktngeneric.Analyzer006,
			testdataDir:    "generic006",
			expectedErrors: 9,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// good.go: 0 errors (proper cmp.Ordered constraint)
			// bad.go: 4 errors (badMin, badMax, badSum, badAverage using ordered ops without constraint)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
