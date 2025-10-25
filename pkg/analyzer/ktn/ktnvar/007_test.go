package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar007(t *testing.T) {
	// 7 slices without capacity (4 original + 3 nouveaux edge cases)
	testhelper.TestGoodBad(t, ktnvar.Analyzer007, "var007", 7)
}
