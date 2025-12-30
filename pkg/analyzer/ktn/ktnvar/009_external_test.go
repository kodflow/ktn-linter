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
			name:           "Large structs passed by value in function params",
			analyzer:       ktnvar.Analyzer009,
			testdataDir:    "var009",
			expectedErrors: 4,
		},
		{
			name:           "Valid pointer usage for large structs",
			analyzer:       ktnvar.Analyzer009,
			testdataDir:    "var009",
			expectedErrors: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 4 grandes structures passées par valeur en paramètres
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
