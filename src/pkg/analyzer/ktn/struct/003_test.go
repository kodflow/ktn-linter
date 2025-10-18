package ktn_struct_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_struct "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/struct"
)

func TestRule003_ExportedFieldsDocumentation(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_struct.Rule003, "struct003")
}
