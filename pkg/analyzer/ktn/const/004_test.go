package ktnconst_test

import (
	"testing"

	ktnconst "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/const"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestConst004(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktnconst.Analyzer004, "const004")
}
