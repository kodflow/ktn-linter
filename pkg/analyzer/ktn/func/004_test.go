package ktnfunc_test

import (
	"testing"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/func"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestFunc004(t *testing.T) {
	// func004/bad.go doit avoir 4 erreurs (naked returns dans fonctions trop longues)
	testhelper.TestGoodBad(t, ktnfunc.Analyzer004, "func004", 4)
}
