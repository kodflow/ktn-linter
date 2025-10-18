package ktn_ops_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_ops "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/ops"
)

func TestOpRule001_DivisionByZero(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_ops.RuleOp001, "op001")
}
