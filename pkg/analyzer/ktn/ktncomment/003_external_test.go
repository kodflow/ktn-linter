// External tests for ktncomment Analyzer003.
package ktncomment_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestComment003 tests the Analyzer003 for constant comment requirement.
//
// Params:
//   - t: testing context
func TestComment003(t *testing.T) {
	// good.go: 0 errors, bad.go: 22 errors for missing constant comments
	testhelper.TestGoodBad(t, ktncomment.Analyzer003, "comment003", 22)
}
