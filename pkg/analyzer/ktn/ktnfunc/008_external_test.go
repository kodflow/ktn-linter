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
	// - Delete: ctx (_ = ctx), req, resp (3 erreurs)
	// - ProcessData: ctx, options (2 erreurs)
	// - PartialIgnore: a, b (_ = b), c (3 erreurs)
	// - BadHandler.Handle: ctx, data (2 erreurs - méthode implémentant interface)
	// Total: 10 erreurs
	testhelper.TestGoodBad(t, ktnfunc.Analyzer008, "func008", 10)
}
