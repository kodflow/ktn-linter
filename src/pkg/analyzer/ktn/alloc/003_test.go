package ktn_alloc_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_alloc "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/alloc"
)

// TestRule003_PreferCompositeLiterals tests the functionality of the corresponding implementation.
func TestRule003_PreferCompositeLiterals(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_alloc.Rule003, "alloc003")
}
