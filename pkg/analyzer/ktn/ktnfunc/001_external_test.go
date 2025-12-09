package ktnfunc_test

import (
	"testing"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestFunc001(t *testing.T) {
	// 6 erreurs de position + 1 erreur "2 types error"
	testhelper.TestGoodBad(t, ktnfunc.Analyzer001, "func001", 7)
}
