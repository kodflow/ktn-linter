package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar007(t *testing.T) {
	// 7 errors: all make([]T, 0) calls without capacity
	// []T{} literals are now in good.go (ignored to avoid false positives)
	testhelper.TestGoodBad(t, ktnvar.Analyzer007, "var007", 7)
}
