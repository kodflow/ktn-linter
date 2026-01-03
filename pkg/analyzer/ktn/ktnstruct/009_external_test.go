// Package ktnstruct_test provides tests for KTN-STRUCT-009.
package ktnstruct_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestAnalyzer009 tests the receiver name consistency analyzer.
func TestAnalyzer009(t *testing.T) {
	testhelper.TestGoodBad(t, ktnstruct.Analyzer009, "struct009", 6)
}
