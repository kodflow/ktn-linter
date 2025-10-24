package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar009(t *testing.T) {
	// 7 maps without capacity hints
	testhelper.TestGoodBad(t, ktnvar.Analyzer009, "var009", 7)
}
