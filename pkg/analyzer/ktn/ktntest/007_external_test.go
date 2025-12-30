package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestTest007 teste l'analyseur KTN-TEST-009.
func TestTest007(t *testing.T) {
	tests := []struct {
		name           string
		dir            string
		expectedErrors int
	}{
		{"good - tests in correct files", "testdata/src/test007/good", 0},
		{"bad - public test in internal error case", "testdata/src/test007/bad", 1},
	}

	// Itération sur les cas de test
	for _, tc := range tests {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			diags := testhelper.RunAnalyzerOnPackage(t, ktntest.Analyzer007, tc.dir)
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
