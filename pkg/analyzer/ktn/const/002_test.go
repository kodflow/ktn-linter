package ktnconst_test

import (
	"testing"

	ktnconst "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/const"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestConst002(t *testing.T) {
	// good.go: 0 errors, bad.go: 6 errors (scattered const blocks + edge cases)
	testhelper.TestGoodBad(t, ktnconst.Analyzer002, "const002", 6)
}

func TestConst002NoVars(t *testing.T) {
	// Test case: const declarations WITHOUT any var declarations
	// good.go: 0 errors (single block), bad.go: 1 error (scattered)
	testhelper.TestGoodBad(t, ktnconst.Analyzer002, "const002_no_vars", 1)
}
