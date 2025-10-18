package ktn_const_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_const "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/const"
)

// TestRule003_IndividualComment tests the functionality of the corresponding implementation.
func TestRule003_IndividualComment(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_const.Rule003, "const003/bad")
	analysistest.Run(t, testdata, ktn_const.Rule003, "const003/good")
}
