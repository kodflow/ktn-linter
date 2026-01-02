package ktngeneric_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktngeneric "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktngeneric"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestGeneric005(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "type parameters shadowing predeclared identifiers",
			analyzer:       ktngeneric.Analyzer005,
			testdataDir:    "generic005",
			expectedErrors: 10,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// good.go: 0 errors (proper type parameter names like T, E, K, V)
			// bad.go: 10 errors (string, int, error, bool, etc. as type params)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
