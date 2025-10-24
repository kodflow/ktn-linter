package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestTest006 tests KTN-TEST-006 rule
func TestTest006(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktntest.Analyzer006, "test006")
}
