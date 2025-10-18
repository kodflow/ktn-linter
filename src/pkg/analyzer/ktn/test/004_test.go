package ktn_test_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_test "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/test"
)

// TestRule004_TestsInTestFiles vérifie que la règle 004 détecte les tests dans des fichiers non-test.
// nolint:KTN-FUNC-001
func TestRule004_TestsInTestFiles(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_test.Rule004, "test004")
}
