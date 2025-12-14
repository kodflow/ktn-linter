// External tests for analyzer 012.
package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestTest012 tests the passthrough test detection.
func TestTest012(t *testing.T) {
	tests := []struct {
		name     string
		analyzer string
	}{
		{
			name:     "passthrough test detection",
			analyzer: "test012",
		},
		{
			name:     "validate test implementation quality",
			analyzer: "test012",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 8 erreurs: 8 tests passthrough dans bad_test.go
			testhelper.TestGoodBadWithFiles(t, ktntest.Analyzer012, tt.analyzer, "good_test.go", "bad_test.go", 8)
		})
	}
}
