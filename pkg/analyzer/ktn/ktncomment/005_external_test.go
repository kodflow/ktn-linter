// External tests for ktncomment Analyzer005.
package ktncomment_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestComment005 tests the Analyzer005 for struct documentation requirement.
//
// Params:
//   - t: testing context
func TestComment005(t *testing.T) {
	// good.go: 0 errors, bad.go: 1 error for missing struct documentation
	testhelper.TestGoodBad(t, ktncomment.Analyzer005, "comment005", 1)
}
