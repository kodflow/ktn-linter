package ktnconst_test

import (
	"testing"

	ktnconst "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/const"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestConst003(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktnconst.Analyzer003, "const003")
}
