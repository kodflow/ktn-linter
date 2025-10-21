package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/var"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar005(t *testing.T) {
	// 13 variables using var with initialization instead of :=
	testhelper.TestGoodBad(t, ktnvar.Analyzer005, "var005", 13)
}
