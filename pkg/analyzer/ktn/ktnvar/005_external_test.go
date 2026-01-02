package ktnvar_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestAnalyzer005 tests the KTN-VAR-005 analyzer.
func TestAnalyzer005(t *testing.T) {
	// Test with expected errors (2 pkg + 1 short decl + 1 block var + 2 range)
	testhelper.TestGoodBad(t, ktnvar.Analyzer005, "var005", 6)
}
