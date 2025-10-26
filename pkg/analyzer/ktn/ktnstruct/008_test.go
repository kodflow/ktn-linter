package ktnstruct_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestStruct008 teste la règle KTN-STRUCT-008.
//
// Returns: aucun
//
// Params:
//   - t: instance de testing
func TestStruct008(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktnstruct.Analyzer008, "struct008")
}
