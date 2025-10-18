package ktnconst_test

import (
	"testing"

	ktnconst "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/const"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestConst005(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktnconst.Analyzer005, "const005")
}
