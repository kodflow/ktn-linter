// External tests for analyzer 013.
package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestTest013 tests the passthrough test detection.
//
// Params:
//   - t: testing context
func TestTest013(t *testing.T) {
	// 3 erreurs: 3 tests passthrough dans bad_test.go
	testhelper.TestGoodBadWithFiles(t, ktntest.Analyzer013, "test013", "good_test.go", "bad_test.go", 3)
}
