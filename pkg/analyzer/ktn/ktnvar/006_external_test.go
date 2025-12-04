package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar006(t *testing.T) {
	// 4 Builder/Buffer declarations without Grow
	testhelper.TestGoodBad(t, ktnvar.Analyzer006, "var006", 4)
}
