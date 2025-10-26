package ktnstruct_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestStruct006 teste la règle KTN-STRUCT-006.
//
// Returns: aucun
//
// Params:
//   - t: instance de testing
func TestStruct006(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktnstruct.Analyzer006, "struct006")
}
