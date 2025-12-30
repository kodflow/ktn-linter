// External tests for analyzer 012.
package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestTest010 tests the passthrough test detection.
func TestTest010(t *testing.T) {
	tests := []struct {
		name     string
		analyzer string
	}{
		{
			name:     "passthrough test detection",
			analyzer: "test010",
		},
		{
			name:     "validate test implementation quality",
			analyzer: "test010",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// 8 erreurs: 8 tests passthrough dans bad_test.go
			testhelper.TestGoodBadWithFiles(t, ktntest.Analyzer010, tt.analyzer, "good_test.go", "bad_test.go", 8)
		})
	}
}
