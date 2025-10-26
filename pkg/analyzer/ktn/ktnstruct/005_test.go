package ktnstruct_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestStruct005 teste la règle KTN-STRUCT-005.
//
// Returns: aucun
//
// Params:
//   - t: instance de testing
func TestStruct005(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktnstruct.Analyzer005, "struct005")
}
