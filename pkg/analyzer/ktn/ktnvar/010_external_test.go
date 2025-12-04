package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar010(t *testing.T) {
	// 8 buffers créés dans des boucles (4 original + 4 nouveaux edge cases)
	testhelper.TestGoodBad(t, ktnvar.Analyzer010, "var010", 8)
}
