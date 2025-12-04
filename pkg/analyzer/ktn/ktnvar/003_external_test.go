package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar003(t *testing.T) {
	// 15 variables using var with initialization instead of := (13 + 2 dans select)
	testhelper.TestGoodBad(t, ktnvar.Analyzer003, "var003", 15)
}
