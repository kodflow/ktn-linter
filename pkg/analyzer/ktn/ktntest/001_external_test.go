package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestAnalyzer001(t *testing.T) {
	// Expected 2 errors in bad package:
	// - helper_test.go (should be renamed to helper_internal_test.go or helper_external_test.go)
	// - resource_test.go (should be renamed to resource_internal_test.go or resource_external_test.go)
	testhelper.TestGoodBadPackage(t, ktntest.Analyzer001, "test001", 2)
}
