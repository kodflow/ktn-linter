package ktnstruct_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestStruct003 teste la règle KTN-STRUCT-003.
//
// Returns: aucun
//
// Params:
//   - t: instance de testing
func TestStruct003(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktnstruct.Analyzer003, "struct003")
}
