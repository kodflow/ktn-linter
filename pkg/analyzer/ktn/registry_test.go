package ktn

import (
	"testing"
)

func TestGetAllRules(t *testing.T) {
	rules := GetAllRules()

	// Check that we have rules
	if len(rules) == 0 {
		t.Error("GetAllRules() returned 0 rules, expected at least 1")
	}

	// Check that all rules are non-nil
	for i, rule := range rules {
		if rule == nil {
			t.Errorf("Rule at index %d is nil", i)
		}
	}

	// Check that we have at least the expected categories (const + func)
	// We should have at least 4 const rules + 4 func rules = 8 rules
	minExpectedRules := 8
	if len(rules) < minExpectedRules {
		t.Errorf("GetAllRules() returned %d rules, expected at least %d", len(rules), minExpectedRules)
	}
}

func TestGetRulesByCategory(t *testing.T) {
	tests := []struct {
		name             string
		category         string
		minExpectedRules int
	}{
		{"const category", "const", 4},
		{"func category", "func", 12},
		{"var category", "var", 19},
		{"test category", "test", 6},
		{"unknown category", "unknown", 0},
		{"empty category", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := GetRulesByCategory(tt.category)

			if tt.minExpectedRules == 0 {
				if len(rules) > 0 {
					t.Errorf("GetRulesByCategory(%q) returned %d rules, expected 0", tt.category, len(rules))
				}
			} else {
				if len(rules) < tt.minExpectedRules {
					t.Errorf("GetRulesByCategory(%q) returned %d rules, expected at least %d", tt.category, len(rules), tt.minExpectedRules)
				}

				// Check that all rules are non-nil
				for i, rule := range rules {
					if rule == nil {
						t.Errorf("Rule at index %d for category %q is nil", i, tt.category)
					}
				}
			}
		})
	}
}
