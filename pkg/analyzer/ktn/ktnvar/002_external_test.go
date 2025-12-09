package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar002(t *testing.T) {
	// 8 variables without explicit type (zero-values are now valid)
	testhelper.TestGoodBad(t, ktnvar.Analyzer002, "var002", 8)
}
