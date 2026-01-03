// Package ktninterface_test provides tests for KTN-INTERFACE-003.
package ktninterface_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktninterface"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestAnalyzer003 tests the single-method interface naming convention analyzer.
func TestAnalyzer003(t *testing.T) {
	testhelper.TestGoodBad(t, ktninterface.Analyzer003, "interface003", 3)
}
