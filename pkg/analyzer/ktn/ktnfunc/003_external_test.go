package ktnfunc_test

import (
	"testing"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestFunc003(t *testing.T) {
	testhelper.TestGoodBad(t, ktnfunc.Analyzer003, "func003", 9)
}
