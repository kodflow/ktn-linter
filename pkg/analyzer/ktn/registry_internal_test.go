package ktn

import (
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestGetAllRules(t *testing.T) {
	const MIN_EXPECTED_RULES int = 8
	tests := []struct {
		name  string
		check func(t *testing.T, rules []*analysis.Analyzer)
	}{
		{
			name: "returns non-empty list",
			check: func(t *testing.T, rules []*analysis.Analyzer) {
				// Vérification liste non vide
				if len(rules) == 0 {
					t.Error("GetAllRules() returned 0 rules, expected at least 1")
				}
			},
		},
		{
			name: "all rules are non-nil",
			check: func(t *testing.T, rules []*analysis.Analyzer) {
				// Vérification règles non-nil
				for i, rule := range rules {
					// Vérification règle
					if rule == nil {
						t.Errorf("Rule at index %d is nil", i)
					}
				}
			},
		},
		{
			name: "has minimum expected rules",
			check: func(t *testing.T, rules []*analysis.Analyzer) {
				// Vérification nombre minimum
				if len(rules) < MIN_EXPECTED_RULES {
					t.Errorf("GetAllRules() returned %d rules, expected at least %d", len(rules), MIN_EXPECTED_RULES)
				}
			},
		},
	}

	rules := GetAllRules()

	// Exécution des tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, rules)
		})
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
		{"var category", "var", 18}, // VAR-010 supprimé
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
