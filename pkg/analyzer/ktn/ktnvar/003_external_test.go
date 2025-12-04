package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar003(t *testing.T) {
	// 6 variables without explicit types + 4 variables with redundant types = 10
	testhelper.TestGoodBad(t, ktnvar.Analyzer003, "var003", 10)
}
