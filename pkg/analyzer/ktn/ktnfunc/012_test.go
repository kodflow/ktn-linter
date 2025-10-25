package ktnfunc_test

import (
	"testing"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestFunc012(t *testing.T) {
	testhelper.TestGoodBad(t, ktnfunc.Analyzer012, "func012", 6)
}
