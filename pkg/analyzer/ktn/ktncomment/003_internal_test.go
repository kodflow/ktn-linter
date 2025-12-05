// Internal tests for 003.go - constant comment analyzer.
package ktncomment

import (
	"testing"
)

// Test_runComment003 tests the runComment003 function configuration.
// The actual analyzer is tested via analysistest in 003_external_test.go.
//
// Params:
//   - t: testing context
func Test_runComment003(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "runComment003 is tested via analysistest"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify that Analyzer003 is properly configured
			if Analyzer003 == nil {
				t.Error("Analyzer003 should not be nil")
				return
			}
			// Check analyzer name
			if Analyzer003.Name != "ktncomment003" {
				t.Errorf("Analyzer003.Name = %q, want %q", Analyzer003.Name, "ktncomment003")
			}
			// Check analyzer requires inspect
			if len(Analyzer003.Requires) == 0 {
				t.Error("Analyzer003 should require inspect.Analyzer")
			}
		})
	}
}
