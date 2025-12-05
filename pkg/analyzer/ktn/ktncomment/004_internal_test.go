// Internal tests for 004.go - variable comment analyzer.
package ktncomment

import (
	"testing"
)

// Test_runComment004 tests the runComment004 function configuration.
// The actual analyzer is tested via analysistest in 004_external_test.go.
//
// Params:
//   - t: testing context
func Test_runComment004(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "runComment004 is tested via analysistest"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify that Analyzer004 is properly configured
			if Analyzer004 == nil {
				t.Error("Analyzer004 should not be nil")
				return
			}
			// Check analyzer name
			if Analyzer004.Name != "ktncomment004" {
				t.Errorf("Analyzer004.Name = %q, want %q", Analyzer004.Name, "ktncomment004")
			}
			// Check analyzer requires inspect
			if len(Analyzer004.Requires) == 0 {
				t.Error("Analyzer004 should require inspect.Analyzer")
			}
		})
	}
}
