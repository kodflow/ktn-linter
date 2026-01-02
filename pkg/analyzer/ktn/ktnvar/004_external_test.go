package ktnvar_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestAnalyzer004 tests the KTN-VAR-004 analyzer.
func TestAnalyzer004(t *testing.T) {
	// Test with expected errors count (2 pkg-level + 2 short decl + 1 block var)
	testhelper.TestGoodBad(t, ktnvar.Analyzer004, "var004", 5)
}
