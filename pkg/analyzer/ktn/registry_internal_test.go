// Internal tests for registry in ktn package.
package ktn

import "testing"

// Test_GetAllRules tests that GetAllRules returns non-empty slice
func Test_GetAllRules(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "returns non-empty slice"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := GetAllRules()
			// Check rules slice is not empty and all rules are non-nil
			if len(rules) == 0 {
				t.Error("GetAllRules() returned empty slice, expected rules")
				return
			}
			for i, rule := range rules {
				if rule == nil {
					t.Errorf("rule at index %d is nil", i)
				}
			}
		})
	}
}

// Test_GetRulesByCategory tests that GetRulesByCategory works correctly
func Test_GetRulesByCategory(t *testing.T) {
	tests := []struct {
		name        string
		category    string
		expectEmpty bool
	}{
		{"const category", "const", false},
		{"func category", "func", false},
		{"struct category", "struct", false},
		{"var category", "var", false},
		{"test category", "test", false},
		{"return category", "return", false},
		{"interface category", "interface", false},
		{"comment category", "comment", false},
		{"modernize category", "modernize", false},
		{"unknown category", "unknown", true},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := GetRulesByCategory(tt.category)

			// Check empty expectation
			if tt.expectEmpty {
				// Should return empty slice for unknown category
				if len(rules) != 0 {
					t.Errorf("expected empty slice for unknown category %q, got %d rules", tt.category, len(rules))
				}
			} else {
				// Should return non-empty slice for known category
				if len(rules) == 0 {
					t.Errorf("expected non-empty slice for category %q", tt.category)
				}

				// Check all rules are non-nil
				for i, rule := range rules {
					// Check rule is not nil
					if rule == nil {
						t.Errorf("rule at index %d is nil for category %q", i, tt.category)
					}
				}
			}
		})
	}
}

// Test_categoryAnalyzers tests that categoryAnalyzers returns valid map
func Test_categoryAnalyzers(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "returns valid map"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categories := categoryAnalyzers()
			// Check that map is not empty and all category functions work
			if len(categories) == 0 {
				t.Error("categoryAnalyzers() returned empty map")
				return
			}
			for name, fn := range categories {
				analyzers := fn()
				if analyzers == nil {
					t.Errorf("category %q returned nil slice", name)
				}
			}
		})
	}
}
