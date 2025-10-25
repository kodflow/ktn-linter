package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar001(t *testing.T) {
	// 6 variables without explicit types
	testhelper.TestGoodBad(t, ktnvar.Analyzer001, "var001", 6)
}
