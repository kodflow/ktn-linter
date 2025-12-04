package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar002(t *testing.T) {
	// good.go: 0 errors (all variables have comments)
	// bad.go: 20 errors (missing comments)
	testhelper.TestGoodBad(t, ktnvar.Analyzer002, "var002", 20)
}
