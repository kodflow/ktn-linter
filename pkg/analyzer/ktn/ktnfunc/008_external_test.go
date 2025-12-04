package ktnfunc_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestFunc008 teste KTN-FUNC-008.
//
// Params:
//   - t: contexte de test
func TestFunc008(t *testing.T) {
	// Expected errors in bad.go:
	// - Delete: req et resp non utilisés (2 erreurs)
	// - ProcessData: ctx et options non utilisés (2 erreurs)
	// - PartialIgnore: a et c non utilisés (2 erreurs)
	// Total: 6 erreurs
	testhelper.TestGoodBad(t, ktnfunc.Analyzer008, "func008", 6)
}
