package ktn_struct_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_struct "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/struct"
)

// TestRule002_GodocDocumentation tests Rule002 godoc documentation requirement.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule002_GodocDocumentation(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_struct.Rule002, "struct002")
}
