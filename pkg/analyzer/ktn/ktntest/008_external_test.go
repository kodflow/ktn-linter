package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestTest008 teste l'analyseur KTN-TEST-010.
func TestTest008(t *testing.T) {
	tests := []struct {
		name           string
		dir            string
		expectedErrors int
	}{
		{"good - tests in correct files", "testdata/src/test008/good", 0},
		{"bad - private test in external error case", "testdata/src/test008/bad", 1},
	}

	// Itération sur les cas de test
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diags := testhelper.RunAnalyzerOnPackage(t, ktntest.Analyzer008, tc.dir)
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
