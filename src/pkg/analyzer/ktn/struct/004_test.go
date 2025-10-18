package ktn_struct_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_struct "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/struct"
)

// TestRule004_MaxFields tests Rule004 maximum fields limit.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule004_MaxFields(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_struct.Rule004, "struct004")
}
