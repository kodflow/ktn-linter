package ktnconst_test

import (
	"testing"

	ktnconst "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestConst004(t *testing.T) {
	// good.go: 0 errors (includes edge cases), bad.go: 33 errors (invalid naming conventions)
	testhelper.TestGoodBad(t, ktnconst.Analyzer004, "const004", 33)
}
