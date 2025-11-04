package ktnvar_test

import (
	"testing"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestVar007(t *testing.T) {
	// 3 errors: only checking make([]T, 0) without capacity
	// []T{} literals are ignored to avoid false positives
	testhelper.TestGoodBad(t, ktnvar.Analyzer007, "var007", 3)
}
