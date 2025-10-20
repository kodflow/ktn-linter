package ktnconst_test

import (
	"testing"

	ktnconst "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/const"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestConst002(t *testing.T) {
	// good.go: 0 errors, bad.go: 6 errors (scattered const blocks + consts after var)
	testhelper.TestGoodBad(t, ktnconst.Analyzer002, "const002", 6)
}
