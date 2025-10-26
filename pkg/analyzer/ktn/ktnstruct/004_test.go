package ktnstruct_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestStruct004 teste la règle KTN-STRUCT-004.
//
// Returns: aucun
//
// Params:
//   - t: instance de testing
func TestStruct004(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktnstruct.Analyzer004, "struct004")
}
