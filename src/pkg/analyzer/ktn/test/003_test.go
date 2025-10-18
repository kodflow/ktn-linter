package ktn_test_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_test "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/test"
)

// TestRule003_NoOrphanTests vérifie que la règle 003 détecte les tests orphelins.
// nolint:KTN-FUNC-001
func TestRule003_NoOrphanTests(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_test.Rule003, "test003")
}

// TestRule003_Extra vérifie que la règle 003 gère les cas supplémentaires correctement.
// nolint:KTN-FUNC-001
func TestRule003_Extra(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_test.Rule003, "test003_extra")
}
