package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar017(t *testing.T) {
	testhelper.TestGoodBad(t, ktnvar.Analyzer017, "var017", 5)
}
