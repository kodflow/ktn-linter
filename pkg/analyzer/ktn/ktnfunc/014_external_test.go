package ktnfunc_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestFunc014 teste KTN-FUNC-014.
//
// Params:
//   - t: contexte de test
func TestFunc014(t *testing.T) {
	// Expected errors in bad.go:
	// - validateTagName: fonction privée non utilisée (créée pour contourner KTN-TEST-008)
	// - unusedHelper: fonction privée non utilisée
	// - formatData: fonction privée non utilisée
	// Total: 3 erreurs
	testhelper.TestGoodBad(t, ktnfunc.Analyzer014, "func014", 3)
}
