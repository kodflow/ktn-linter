package ktn_ops_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_ops "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/ops"
)

// TestReturnRule001_NakedReturns tests the functionality of the corresponding implementation.
func TestReturnRule001_NakedReturns(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_ops.RuleReturn001, "return001")
}
