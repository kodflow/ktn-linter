package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar002(t *testing.T) {
	// 8 variables without explicit type + 2 variables without value = 10
	testhelper.TestGoodBad(t, ktnvar.Analyzer002, "var002", 10)
}
