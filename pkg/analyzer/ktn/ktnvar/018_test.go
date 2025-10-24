package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar018 vérifie la détection des conversions string() répétées.
//
// Params:
//   - t: contexte de test
func TestVar018(t *testing.T) {
	// 12 conversions répétées détectées (5 original + 7 nouveaux edge cases)
	testhelper.TestGoodBad(t, ktnvar.Analyzer018, "var018", 12)
}
