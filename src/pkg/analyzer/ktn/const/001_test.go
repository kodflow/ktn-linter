package ktn_const_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_const "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/const"
)

func TestRule001_ConstGrouping(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_const.Rule001, "const001")
}
