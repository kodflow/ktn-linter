package ktngeneric_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktngeneric "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktngeneric"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestGeneric002(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "unnecessary generics on interface types",
			analyzer:       ktngeneric.Analyzer002,
			testdataDir:    "generic002",
			expectedErrors: 4,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// good.go: 0 errors (generics used correctly or interfaces used directly)
			// bad.go: 3 errors (unnecessary generics on interface constraints)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
