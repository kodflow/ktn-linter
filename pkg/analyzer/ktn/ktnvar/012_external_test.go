package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar012 vérifie la détection des conversions string() répétées.
func TestVar012(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Repeated string() conversions",
			analyzer:       ktnvar.Analyzer012,
			testdataDir:    "var012",
			expectedErrors: 12,
		},
		{
			name:           "Valid cached string conversions",
			analyzer:       ktnvar.Analyzer012,
			testdataDir:    "var012",
			expectedErrors: 12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 12 conversions répétées détectées (5 original + 7 nouveaux edge cases)
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
