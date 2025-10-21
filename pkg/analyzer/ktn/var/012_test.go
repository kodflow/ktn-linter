package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/var"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar012(t *testing.T) {
	// 6 string concatenation errors detected
	testhelper.TestGoodBad(t, ktnvar.Analyzer012, "var012", 6)
}
