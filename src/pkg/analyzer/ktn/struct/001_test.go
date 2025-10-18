package ktn_struct_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_struct "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/struct"
)

// TestRule001_MixedCapsNaming tests Rule001 MixedCaps naming convention.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule001_MixedCapsNaming(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_struct.Rule001, "struct001")
}
