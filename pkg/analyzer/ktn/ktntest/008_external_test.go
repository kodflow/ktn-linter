package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestTest008 teste l'analyseur KTN-TEST-008 (règle 1:2).
//
// Params:
//   - t: instance de test
func TestTest008(t *testing.T) {
	// Test du package good (doit avoir 0 erreur - les deux fichiers de test sont présents)
	goodDir := "testdata/src/test008/good"
	diags := testhelper.RunAnalyzerOnPackage(t, ktntest.Analyzer008, goodDir)
	if len(diags) != 0 {
		t.Errorf("%s should have 0 errors, got %d", goodDir, len(diags))
		for _, d := range diags {
			t.Logf("  %v", d.Message)
		}
	}

	// Test bad_no_tests (doit avoir 1 erreur - aucun fichier de test)
	badNoTestsDir := "testdata/src/test008/bad_no_tests"
	diags = testhelper.RunAnalyzerOnPackage(t, ktntest.Analyzer008, badNoTestsDir)
	if len(diags) != 1 {
		t.Errorf("%s should have 1 error, got %d", badNoTestsDir, len(diags))
		for _, d := range diags {
			t.Logf("  %v", d.Message)
		}
	}

	// Test bad_only_internal (doit avoir 1 erreur - manque _external_test.go)
	badOnlyInternalDir := "testdata/src/test008/bad_only_internal"
	diags = testhelper.RunAnalyzerOnPackage(t, ktntest.Analyzer008, badOnlyInternalDir)
	if len(diags) != 1 {
		t.Errorf("%s should have 1 error, got %d", badOnlyInternalDir, len(diags))
		for _, d := range diags {
			t.Logf("  %v", d.Message)
		}
	}

	// Test bad_only_external (doit avoir 1 erreur - manque _internal_test.go)
	badOnlyExternalDir := "testdata/src/test008/bad_only_external"
	diags = testhelper.RunAnalyzerOnPackage(t, ktntest.Analyzer008, badOnlyExternalDir)
	if len(diags) != 1 {
		t.Errorf("%s should have 1 error, got %d", badOnlyExternalDir, len(diags))
		for _, d := range diags {
			t.Logf("  %v", d.Message)
		}
	}
}
