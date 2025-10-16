package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
)

// TestErrorAnalyzer teste l'analyseur error (KTN-ERROR-001).
//
// Params:
//   - t: l'instance de test
func TestErrorAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.ErrorAnalyzer, "error")
}
