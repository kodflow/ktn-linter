package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/var"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar019 vérifie la détection des copies de mutex.
//
// Params:
//   - t: contexte de test
func TestVar019(t *testing.T) {
	// 9 cas de copies de mutex attendus
	testhelper.TestGoodBad(t, ktnvar.Analyzer019, "var019", 9)
}
