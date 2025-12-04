package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar011 vérifie la détection du shadowing de variables.
//
// Params:
//   - t: contexte de test
func TestVar011(t *testing.T) {
	// 5 cas de shadowing attendus
	testhelper.TestGoodBad(t, ktnvar.Analyzer011, "var011", 5)
}
