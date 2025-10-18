package ktn_data_structures_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_data_structures "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/data_structures"
)

// TestSliceRule001_IndexWithoutBoundsCheck tests the functionality of the corresponding implementation.
func TestSliceRule001_IndexWithoutBoundsCheck(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_data_structures.RuleSlice001, "slice001")
}
