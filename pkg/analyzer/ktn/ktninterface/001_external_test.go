// External tests for ktninterface Analyzer001.
package ktninterface_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktninterface"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestInterface001 tests the Analyzer001 for unused interface declarations.
func TestInterface001(t *testing.T) {
	// good.go: 0 errors, bad.go: 3 errors for unused interfaces
	testhelper.TestGoodBad(t, ktninterface.Analyzer001, "interface001", 3)
}
