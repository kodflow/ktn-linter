// External tests for ktncomment Analyzer006.
package ktncomment_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestComment006 tests the Analyzer006 for function documentation format.
//
// Params:
//   - t: testing context
func TestComment006(t *testing.T) {
	// good.go: 0 errors, bad.go: 8 errors for missing/invalid function docs
	testhelper.TestGoodBad(t, ktncomment.Analyzer006, "comment006", 8)
}
