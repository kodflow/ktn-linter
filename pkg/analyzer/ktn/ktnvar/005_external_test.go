package ktnvar_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestAnalyzer005 tests the KTN-VAR-005 analyzer.
func TestAnalyzer005(t *testing.T) {
	// Test with expected errors count
	testhelper.TestGoodBad(t, ktnvar.Analyzer005, "var005", 2)
}
