// Package ktnstruct_test provides tests for KTN-STRUCT-008.
package ktnstruct_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestAnalyzer008 tests the receiver type consistency analyzer.
func TestAnalyzer008(t *testing.T) {
	testhelper.TestGoodBad(t, ktnstruct.Analyzer008, "struct008", 3)
}
