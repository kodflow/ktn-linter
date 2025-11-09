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
	// Test du package good (doit avoir 0 erreur - les tests sont dans les bons fichiers)
	goodDir := "testdata/src/test010/good"
	diags := testhelper.RunAnalyzerOnPackage(t, ktntest.Analyzer010, goodDir)
	if len(diags) != 0 {
		t.Errorf("%s should have 0 errors, got %d", goodDir, len(diags))
		for _, d := range diags {
			t.Logf("  %v", d.Message)
		}
	}

	// Test bad (doit avoir 1 erreur - test de fonction privée dans _external_test.go)
	badDir := "testdata/src/test010/bad"
	diags = testhelper.RunAnalyzerOnPackage(t, ktntest.Analyzer010, badDir)
	if len(diags) != 1 {
		t.Errorf("%s should have 1 error, got %d", badDir, len(diags))
		for _, d := range diags {
			t.Logf("  %v", d.Message)
		}
	}
}
