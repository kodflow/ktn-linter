package ktnfunc_test

import (
	"testing"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestFunc012(t *testing.T) {
	// Placeholder: will calculate error count by running test
	testhelper.TestGoodBad(t, ktnfunc.Analyzer012, "func012", 4)
}
