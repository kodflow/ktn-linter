package ktnconst_test

import (
	"testing"

	ktnconst "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestConst002(t *testing.T) {
	// good.go: 0 errors
	// bad.go: 10 errors (4 scattered + 3 after var + 2 after type + 1 after func)
	testhelper.TestGoodBad(t, ktnconst.Analyzer002, "const002", 10)
}
