package ktnfunc_test

import (
	"testing"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/func"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestFunc005(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktnfunc.Analyzer005, "func005")
}
