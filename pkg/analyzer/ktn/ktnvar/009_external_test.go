package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar009(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Large structs without pointer",
			analyzer:       ktnvar.Analyzer009,
			testdataDir:    "var009",
			expectedErrors: 5,
		},
		{
			name:           "Valid pointer usage for large structs",
			analyzer:       ktnvar.Analyzer009,
			testdataDir:    "var009",
			expectedErrors: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 5 grandes structures sans pointeur
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
