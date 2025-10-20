package ktnfunc_test

import (
	"testing"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/func"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestFunc003(t *testing.T) {
	// Placeholder: will calculate error count by running test
	testhelper.TestGoodBad(t, ktnfunc.Analyzer003, "func003", 12)
}


