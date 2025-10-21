package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/var"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar008(t *testing.T) {
	// 3 make calls with length > 0
	testhelper.TestGoodBad(t, ktnvar.Analyzer008, "var008", 3)
}
