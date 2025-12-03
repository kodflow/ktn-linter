// External tests for ktnreturn Analyzer002.
package ktnreturn_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnreturn"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestReturn002 tests the Analyzer002 for nil returns in slice/map types.
func TestReturn002(t *testing.T) {
	// good.go: 0 errors, bad.go: 6 errors for nil slice/map returns
	testhelper.TestGoodBad(t, ktnreturn.Analyzer002, "return002", 6)
}
