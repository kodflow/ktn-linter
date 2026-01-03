// Package ktninterface_test provides tests for KTN-INTERFACE-004.
package ktninterface_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktninterface"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestAnalyzer004 tests the empty interface overuse analyzer.
func TestAnalyzer004(t *testing.T) {
	testhelper.TestGoodBad(t, ktninterface.Analyzer004, "interface004", 5)
}
