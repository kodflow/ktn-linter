package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar006 vérifie que les variables de package sont déclarées après les constantes.
//
// Params:
//   - t: instance de test
func TestVar006(t *testing.T) {
	testhelper.TestGoodBad(t, ktnvar.Analyzer006, "var006", 1)
}
