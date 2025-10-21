package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/var"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar004(t *testing.T) {
	// good.go: 0 errors (all variables have comments)
	// bad.go: 20 errors (missing comments)
	testhelper.TestGoodBad(t, ktnvar.Analyzer004, "var004", 20)
}
