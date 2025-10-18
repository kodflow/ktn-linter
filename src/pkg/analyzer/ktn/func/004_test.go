package ktn_func_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_func "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/func"
)

// TestRule004_ReturnsDocumentation tests the functionality of the corresponding implementation.
func TestRule004_ReturnsDocumentation(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_func.Rule004, "func004")
}
