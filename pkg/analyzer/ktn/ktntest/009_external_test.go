package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestTest009 teste l'analyseur KTN-TEST-011.
func TestTest009(t *testing.T) {
	tests := []struct {
		name           string
		dir            string
		expectedErrors int
	}{
		{"good - conventions respected", "testdata/src/test009/good", 0},
		{"bad - internal with wrong package error case", "testdata/src/test009/bad_internal", 1},
		{"bad - external with wrong package error case", "testdata/src/test009/bad_external", 1},
	}

	// Itération sur les cas de test
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diags := testhelper.RunAnalyzerOnPackage(t, ktntest.Analyzer009, tc.dir)
			// Vérification du nombre de diagnostics
			if len(diags) != tc.expectedErrors {
				t.Errorf("%s should have %d errors, got %d", tc.dir, tc.expectedErrors, len(diags))
				// Affichage des diagnostics
				for _, d := range diags {
					t.Logf("  %v", d.Message)
				}
			}
		})
	}
}
