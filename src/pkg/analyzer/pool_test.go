package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
)

// TestPoolAnalyzer teste l'analyseur pool (KTN-POOL-001 et KTN-STRUCT-004).
//
// Params:
//   - t: l'instance de test
func TestPoolAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.PoolAnalyzer, "pool")
}
