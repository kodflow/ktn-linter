package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/var"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar003(t *testing.T) {
	// 6 variables with SCREAMING_SNAKE_CASE naming
	testhelper.TestGoodBad(t, ktnvar.Analyzer003, "var003", 6)
}
