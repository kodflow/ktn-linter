package ktn_const_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_const "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/const"
)

// TestRule004_ExplicitType tests the functionality of the corresponding implementation.
func TestRule004_ExplicitType(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_const.Rule004, "const004/bad", "const004/good")
}
