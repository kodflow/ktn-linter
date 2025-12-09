package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar018(t *testing.T) {
	// 8 variables using snake_case (with underscores)
	testhelper.TestGoodBad(t, ktnvar.Analyzer018, "var018", 8)
}
