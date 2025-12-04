package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar014 vérifie que les variables de package sont déclarées après les constantes.
//
// Params:
//   - t: instance de test
func TestVar014(t *testing.T) {
	testhelper.TestGoodBad(t, ktnvar.Analyzer014, "var014", 1)
}
