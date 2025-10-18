package ktn_test_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_test "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/test"
)

func TestRule002_CoverageRequired(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_test.Rule002, "test002")
}
