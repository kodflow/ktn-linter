package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar005(t *testing.T) {
	// 8 make calls with length > 0 (VAR-016 cases excluded)
	testhelper.TestGoodBad(t, ktnvar.Analyzer005, "var005", 8)
}
