package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar006(t *testing.T) {
	// 9 make calls with length > 0 (3 original + 6 edge cases)
	testhelper.TestGoodBad(t, ktnvar.Analyzer006, "var006", 9)
}
