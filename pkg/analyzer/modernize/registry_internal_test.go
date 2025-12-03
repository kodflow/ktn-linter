// Internal tests for registry in modernize package.
package modernize

import "testing"

// Test_Analyzers tests that Analyzers returns non-empty slice
func Test_Analyzers(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"verify analyzers returned correctly"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzers := Analyzers()

			// Check that analyzers slice is not empty
			if len(analyzers) == 0 {
				t.Error("Analyzers() returned empty slice, expected analyzers")
			}

			// Check that all analyzers are non-nil
			for i, analyzer := range analyzers {
				// Check analyzer is not nil
				if analyzer == nil {
					t.Errorf("analyzer at index %d is nil", i)
				}
			}

			// Verify that disabled analyzers are not present
			disabledNames := []string{"newexpr"}
			for _, analyzer := range analyzers {
				// Check that no disabled analyzer is present
				for _, disabled := range disabledNames {
					// Check analyzer is not in disabled list
					if analyzer.Name == disabled {
						t.Errorf("disabled analyzer %q should not be in the list", disabled)
					}
				}
			}
		})
	}
}
