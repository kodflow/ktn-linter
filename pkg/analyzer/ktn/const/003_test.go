package ktnconst_test

import (
	"testing"

	ktnconst "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/const"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestConst003(t *testing.T) {
	// good.go: 0 errors (includes edge cases), bad.go: 33 errors (invalid naming conventions)
	testhelper.TestGoodBad(t, ktnconst.Analyzer003, "const003", 33)
}
