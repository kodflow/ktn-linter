// Package ktnstruct_test provides tests for KTN-STRUCT-007.
package ktnstruct_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestAnalyzer007 tests the DTO exported fields serialization tag analyzer.
func TestAnalyzer007(t *testing.T) {
	testhelper.TestGoodBad(t, ktnstruct.Analyzer007, "struct007", 6)
}
