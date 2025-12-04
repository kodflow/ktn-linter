package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar007(t *testing.T) {
	// 4 Builder/Buffer declarations without Grow
	testhelper.TestGoodBad(t, ktnvar.Analyzer007, "var007", 4)
}
