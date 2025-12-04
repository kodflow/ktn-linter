package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar004(t *testing.T) {
	// 15 variables using var with initialization instead of := (13 + 2 dans select)
	testhelper.TestGoodBad(t, ktnvar.Analyzer004, "var004", 15)
}
