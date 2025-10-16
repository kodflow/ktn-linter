package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
)

// TestGoroutineAnalyzer teste l'analyseur goroutine (KTN-GOROUTINE-001 et KTN-GOROUTINE-002).
//
// Params:
//   - t: l'instance de test
func TestGoroutineAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.GoroutineAnalyzer, "goroutine")
}
