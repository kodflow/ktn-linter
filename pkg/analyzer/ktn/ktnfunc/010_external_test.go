package ktnfunc_test

import (
	"testing"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestFunc010(t *testing.T) {
	// func010/bad.go doit avoir 4 erreurs (naked returns dans fonctions trop longues)
	testhelper.TestGoodBad(t, ktnfunc.Analyzer010, "func010", 4)
}
