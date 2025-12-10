// Internal tests for registry in ktn package.
package ktn

import "testing"

// Test_categoryAnalyzers tests that categoryAnalyzers returns valid map.
//
// Params:
//   - t: testing context
func Test_categoryAnalyzers(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "returns valid map"},
	}

	// Iteration over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categories := categoryAnalyzers()
			// Check that map is not empty
			if len(categories) == 0 {
				t.Error("categoryAnalyzers() returned empty map")
				return
			}
			// Check all category functions work
			for name, fn := range categories {
				analyzers := fn()
				// Check analyzers are not nil
				if analyzers == nil {
					t.Errorf("category %q returned nil slice", name)
				}
			}
		})
	}
}
