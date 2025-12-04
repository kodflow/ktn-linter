// External tests for ktncomment Analyzer002.
package ktncomment_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestComment002 tests the Analyzer002 for package comment requirement.
//
// Params:
//   - t: testing context
func TestComment002(t *testing.T) {
	// good/: 0 errors, bad/: 2 errors for missing package comments
	testhelper.TestGoodBadPackage(t, ktncomment.Analyzer002, "comment002", 2)
}
