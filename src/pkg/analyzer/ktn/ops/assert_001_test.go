package ktn_ops_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_ops "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/ops"
)

// TestAssertRule001_TypeAssertionWithoutOk tests the functionality of the corresponding implementation.
func TestAssertRule001_TypeAssertionWithoutOk(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_ops.RuleAssert001, "assert001")
}
