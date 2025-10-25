package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar011(t *testing.T) {
	// 4 Builder/Buffer declarations without Grow
	testhelper.TestGoodBad(t, ktnvar.Analyzer011, "var011", 4)
}
