package ktn_ops_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_ops "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/ops"
)

func TestCompRule001_RedundantComparisons(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_ops.RuleComp001, "comp001")
}
