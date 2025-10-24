package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar013(t *testing.T) {
	// 7 allocations dans des boucles (5 assignements + 2 déclarations var)
	testhelper.TestGoodBad(t, ktnvar.Analyzer013, "var013", 7)
}
