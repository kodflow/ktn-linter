package ktn_test_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_test "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/test"
)

// TestRule001_PackageTestSuffix vérifie que la règle 001 détecte les packages de test sans suffixe _test.
// nolint:KTN-FUNC-001
func TestRule001_PackageTestSuffix(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_test.Rule001, "test001")
}
