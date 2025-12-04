// External tests for ktncomment Analyzer007.
package ktncomment_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestComment007 tests the Analyzer007 for control block comments.
//
// Params:
//   - t: testing context
func TestComment007(t *testing.T) {
	// good.go: 0 errors, bad.go: 33 errors for missing block comments
	testhelper.TestGoodBad(t, ktncomment.Analyzer007, "comment007", 33)
}
