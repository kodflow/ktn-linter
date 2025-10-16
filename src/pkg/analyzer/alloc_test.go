package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
)

// TestAllocAnalyzer teste l'analyseur d'allocation.
//
// Params:
//   - t: l'instance de test
func TestAllocAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.AllocAnalyzer, "alloc")
}
