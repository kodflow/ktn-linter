package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar018 vérifie la détection des copies de mutex.
//
// Params:
//   - t: contexte de test
func TestVar018(t *testing.T) {
	// 15 cas de copies de mutex attendus
	testhelper.TestGoodBad(t, ktnvar.Analyzer018, "var018", 15)
}
