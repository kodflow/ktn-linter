package ktnconst_test

import (
	"testing"

	ktnconst "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestConst002(t *testing.T) {
	// good.go: 0 errors, bad.go: 7 errors (6 scattered + 1 after var)
	testhelper.TestGoodBad(t, ktnconst.Analyzer002, "const002", 7)
}
