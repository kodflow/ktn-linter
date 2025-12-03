package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestTest010 teste l'analyseur KTN-TEST-010.
//
// Params:
//   - t: instance de test
func TestTest010(t *testing.T) {
	tests := []struct {
		name           string
		dir            string
		expectedErrors int
	}{
		{"good - tests in correct files", "testdata/src/test010/good", 0},
		{"bad - private test in external error case", "testdata/src/test010/bad", 1},
	}

	// Itération sur les cas de test
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diags := testhelper.RunAnalyzerOnPackage(t, ktntest.Analyzer010, tc.dir)
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
