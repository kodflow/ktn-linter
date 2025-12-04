package ktnconst_test

import (
	"testing"

	ktnconst "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestConst001(t *testing.T) {
	// good.go: 0 errors, bad.go: 13 errors (constants without explicit type)
	// - 10 in block + 1 iota + 2 multi-name
	testhelper.TestGoodBad(t, ktnconst.Analyzer001, "const001", 13)
}
