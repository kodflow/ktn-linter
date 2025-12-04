package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar005(t *testing.T) {
	// 8 errors:
	// - 7 make([]T, 0) calls without capacity
	// - 1 []T{} literal followed by append (should use make)
	// []T{} in returns/structs are NOT reported (avoid false positives)
	testhelper.TestGoodBad(t, ktnvar.Analyzer005, "var005", 8)
}
