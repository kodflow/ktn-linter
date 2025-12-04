package ktnconst_test

import (
	"testing"

	ktnconst "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestConst003(t *testing.T) {
	// good.go: 0 errors, bad.go: 8 errors (7 scattered + 1 after var)
	testhelper.TestGoodBad(t, ktnconst.Analyzer003, "const003", 8)
}
