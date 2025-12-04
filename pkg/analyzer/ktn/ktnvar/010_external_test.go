package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar010(t *testing.T) {
	// 5 grandes structures sans pointeur
	testhelper.TestGoodBad(t, ktnvar.Analyzer010, "var010", 5)
}
