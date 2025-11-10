package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestTest011 teste l'analyseur KTN-TEST-011.
//
// Params:
//   - t: instance de test
func TestTest011(t *testing.T) {
	// Test du package good (doit avoir 0 erreur - conventions respectées)
	goodDir := "testdata/src/test011/good"
	diags := testhelper.RunAnalyzerOnPackage(t, ktntest.Analyzer011, goodDir)
	if len(diags) != 0 {
		t.Errorf("%s should have 0 errors, got %d", goodDir, len(diags))
		for _, d := range diags {
			t.Logf("  %v", d.Message)
		}
	}

	// Test bad_internal (doit avoir 1 erreur - _internal_test.go avec package xxx_test)
	badInternalDir := "testdata/src/test011/bad_internal"
	diags = testhelper.RunAnalyzerOnPackage(t, ktntest.Analyzer011, badInternalDir)
	if len(diags) != 1 {
		t.Errorf("%s should have 1 error, got %d", badInternalDir, len(diags))
		for _, d := range diags {
			t.Logf("  %v", d.Message)
		}
	}

	// Test bad_external (doit avoir 1 erreur - _external_test.go avec package xxx)
	badExternalDir := "testdata/src/test011/bad_external"
	diags = testhelper.RunAnalyzerOnPackage(t, ktntest.Analyzer011, badExternalDir)
	if len(diags) != 1 {
		t.Errorf("%s should have 1 error, got %d", badExternalDir, len(diags))
		for _, d := range diags {
			t.Logf("  %v", d.Message)
		}
	}
}
