package ktnfunc_test

import (
	"testing"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestFunc012(t *testing.T) {
	// func012/bad.go doit avoir 4 erreurs (naked returns dans fonctions trop longues)
	testhelper.TestGoodBad(t, ktnfunc.Analyzer012, "func012", 4)
}
