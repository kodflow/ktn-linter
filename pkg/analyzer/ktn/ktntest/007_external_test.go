package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestTest007 teste la règle KTN-TEST-007
func TestTest007(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktntest.Analyzer007, "test007")
}
