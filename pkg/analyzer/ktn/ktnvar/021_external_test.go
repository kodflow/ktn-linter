package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar021 vérifie la détection des incohérences de receiver.
// Erreurs attendues dans bad.go:
// - Server.Stop: valeur au lieu de pointeur (1)
// - Server.Status: valeur au lieu de pointeur (1)
// Total: 2 erreurs
func TestVar021(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Inconsistent receiver types",
			analyzer:       ktnvar.Analyzer021,
			testdataDir:    "var021",
			expectedErrors: 2,
		},
		{
			name:           "Consistent receiver types",
			analyzer:       ktnvar.Analyzer021,
			testdataDir:    "var021",
			expectedErrors: 2,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
