// External tests for analyzer 012.
package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestTest012 tests the passthrough test detection.
//
// Params:
//   - t: testing context
func TestTest012(t *testing.T) {
	// 8 erreurs: 8 tests passthrough dans bad_test.go
	testhelper.TestGoodBadWithFiles(t, ktntest.Analyzer012, "test012", "good_test.go", "bad_test.go", 8)
}
