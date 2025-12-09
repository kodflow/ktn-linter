package ktnconst_test

import (
	"testing"

	ktnconst "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestConst003(t *testing.T) {
	// good.go: 0 errors (CamelCase is valid)
	// bad.go: 15 errors (constants with underscores violate CamelCase convention)
	testhelper.TestGoodBad(t, ktnconst.Analyzer003, "const003", 15)
}
