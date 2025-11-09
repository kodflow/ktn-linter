package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestTest008 teste l'analyseur KTN-TEST-008.
//
// Params:
//   - t: instance de test
func TestTest008(t *testing.T) {
	// test008/ doit avoir 3 erreurs dans le package bad :
	// - code_test.go : manque suffixe _internal ou _external
	// - wrong_internal_test.go : mauvais package (devrait être test008)
	// - wrong_external_test.go : mauvais package (devrait être test008_test)
	testhelper.TestGoodBadPackage(t, ktntest.Analyzer008, "test008", 3)
}
