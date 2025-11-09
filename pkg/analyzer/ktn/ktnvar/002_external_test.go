package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar002(t *testing.T) {
	// 5 scattered var declarations (2 single + 3 groups after first)
	testhelper.TestGoodBad(t, ktnvar.Analyzer002, "var002", 4)
}
