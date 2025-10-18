package ktn_func_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_func "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/func"
)

func TestRule008_ReturnComments(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_func.Rule008, "func008")
}
