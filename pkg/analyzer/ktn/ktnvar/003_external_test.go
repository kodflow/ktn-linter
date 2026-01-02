package ktnvar_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestAnalyzer003 tests the KTN-VAR-003 analyzer.
func TestAnalyzer003(t *testing.T) {
	// Test with expected errors count
	testhelper.TestGoodBad(t, ktnvar.Analyzer003, "var003", 3)
}
