package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/var"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar010(t *testing.T) {
	// 8 empty slice literals
	testhelper.TestGoodBad(t, ktnvar.Analyzer010, "var010", 8)
}
