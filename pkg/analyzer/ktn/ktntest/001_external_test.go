package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestAnalyzer001(t *testing.T) {
	tests := []struct {
		name     string
		analyzer string
	}{
		{
			name:     "test file naming validation",
			analyzer: "test001",
		},
		{
			name:     "verify suffix compliance",
			analyzer: "test001",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Expected 2 errors in bad package:
			// - helper_test.go (should be renamed to helper_internal_test.go or helper_external_test.go)
			// - resource_test.go (should be renamed to resource_internal_test.go or resource_external_test.go)
			testhelper.TestGoodBadPackage(t, ktntest.Analyzer001, tt.analyzer, 2)
		})
	}
}
