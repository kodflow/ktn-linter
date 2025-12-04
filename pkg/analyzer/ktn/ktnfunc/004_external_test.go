package ktnfunc_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestFunc004 teste KTN-FUNC-004.
//
// Params:
//   - t: contexte de test
func TestFunc004(t *testing.T) {
	// Expected errors in bad.go:
	// - validateTagName: fonction privée non utilisée (créée pour contourner KTN-TEST-008)
	// - unusedHelper: fonction privée non utilisée
	// - formatData: fonction privée non utilisée
	// Total: 3 erreurs
	testhelper.TestGoodBad(t, ktnfunc.Analyzer004, "func004", 3)
}

// TestFunc014_SpecialFunctions teste que main, init et callbacks sont ignorés.
//
// Params:
//   - t: contexte de test
func TestFunc014_SpecialFunctions(t *testing.T) {
	// Expected errors in bad.go:
	// - deadFunction: fonction privée non utilisée
	// Total: 1 erreur
	//
	// good.go ne doit pas générer d'erreur:
	// - main() est ignorée (point d'entrée)
	// - init() est ignorée (appelée automatiquement)
	// - run() est utilisée comme callback (RunE: run)
	// - helper() est assignée à une variable (var helperFunc = helper)
	testhelper.TestGoodBad(t, ktnfunc.Analyzer004, "func004_special", 1)
}
