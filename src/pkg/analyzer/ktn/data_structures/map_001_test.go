package ktn_data_structures_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_data_structures "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/data_structures"
)

// TestMapRule001_WriteWithoutNilCheck tests the functionality of the corresponding implementation.
func TestMapRule001_WriteWithoutNilCheck(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_data_structures.RuleMap001, "map001")
}
