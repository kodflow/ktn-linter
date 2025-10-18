package ktn_const_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_const "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/const"
)

// TestRule002_GroupComment tests the functionality of the corresponding implementation.
func TestRule002_GroupComment(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_const.Rule002, "const002/bad", "const002/good")
}
