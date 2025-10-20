package ktnconst_test

import (
	"testing"

	ktnconst "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/const"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestConst001(t *testing.T) {
	// good.go: 0 errors, bad.go: 10 errors (one per constant without explicit type)
	testhelper.TestGoodBad(t, ktnconst.Analyzer001, "const001", 10)
}
