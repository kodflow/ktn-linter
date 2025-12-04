package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar008(t *testing.T) {
	// 6 string concatenation errors detected
	testhelper.TestGoodBad(t, ktnvar.Analyzer008, "var008", 6)
}
