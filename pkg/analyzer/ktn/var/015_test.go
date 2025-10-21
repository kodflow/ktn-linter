package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/var"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar015(t *testing.T) {
	testhelper.TestGoodBad(t, ktnvar.Analyzer015, "var015", 4)
}
