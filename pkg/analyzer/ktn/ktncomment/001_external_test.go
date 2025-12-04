// External tests for ktncomment Analyzer001.
package ktncomment_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestComment001 tests the Analyzer001 for inline comments exceeding 80 characters.
func TestComment001(t *testing.T) {
	// good.go: 0 errors, bad.go: 4 errors for long inline comments
	testhelper.TestGoodBad(t, ktncomment.Analyzer001, "comment001", 4)
}
