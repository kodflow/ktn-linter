package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/var"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar017 vérifie la détection du shadowing de variables.
//
// Params:
//   - t: contexte de test
func TestVar017(t *testing.T) {
	// 5 cas de shadowing attendus
	testhelper.TestGoodBad(t, ktnvar.Analyzer017, "var017", 5)
}
