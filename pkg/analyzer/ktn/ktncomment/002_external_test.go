// External tests for ktncomment Analyzer002.
package ktncomment_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestComment002 tests the Analyzer002 for inline comments exceeding 80 characters.
func TestComment002(t *testing.T) {
	// good.go: 0 errors, bad.go: 4 errors for long inline comments
	testhelper.TestGoodBad(t, ktncomment.Analyzer002, "comment002", 4)
}
