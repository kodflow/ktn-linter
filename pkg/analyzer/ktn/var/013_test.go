package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/var"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar013(t *testing.T) {
	// 5 allocations dans des boucles
	testhelper.TestGoodBad(t, ktnvar.Analyzer013, "var013", 5)
}
