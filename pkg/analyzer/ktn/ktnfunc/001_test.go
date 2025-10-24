package ktnfunc_test

import (
	"testing"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestFunc001(t *testing.T) {
	// func001/bad.go doit avoir 2 erreurs (TooLong: 36 lignes, VeryLong: 37 lignes)
	testhelper.TestGoodBad(t, ktnfunc.Analyzer001, "func001", 2)
}

