package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar016(t *testing.T) {
	testhelper.TestGoodBad(t, ktnvar.Analyzer016, "var016", 5)
}
