package ktnconst_test

import (
	"testing"

	ktnconst "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestConst002(t *testing.T) {
	// good.go: 0 errors (includes edge cases), bad.go: 22 errors (missing comments + want directives)
	testhelper.TestGoodBad(t, ktnconst.Analyzer002, "const002", 22)
}
