package ktn_ops_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_ops "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/ops"
)

// TestPointerRule001_NilDereference tests the functionality of the corresponding implementation.
func TestPointerRule001_NilDereference(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_ops.RulePointer001, "pointer001")
}
