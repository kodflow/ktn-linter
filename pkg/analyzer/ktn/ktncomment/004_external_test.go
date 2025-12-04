// External tests for ktncomment Analyzer004.
package ktncomment_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestComment004 tests the Analyzer004 for variable comment requirement.
//
// Params:
//   - t: testing context
func TestComment004(t *testing.T) {
	// good.go: 0 errors, bad.go: 20 errors for missing variable comments
	testhelper.TestGoodBad(t, ktncomment.Analyzer004, "comment004", 20)
}
