package ktnfunc_test

import (
	"testing"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestFunc011(t *testing.T) {
	testhelper.TestGoodBad(t, ktnfunc.Analyzer011, "func011", 9)
}
