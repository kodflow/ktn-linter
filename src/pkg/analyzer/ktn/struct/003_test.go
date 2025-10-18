package ktn_struct_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_struct "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/struct"
)

// TestRule003_ExportedFieldsDocumentation tests Rule003 exported fields documentation requirement.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule003_ExportedFieldsDocumentation(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_struct.Rule003, "struct003")
}
