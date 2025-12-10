package ktn_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn"
	"golang.org/x/tools/go/analysis"
)

// TestGetAllRules tests the GetAllRules function with error cases
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

	rules := ktn.GetAllRules()

	// Exécution des tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, rules)
		})
	}
}

// TestGetRulesByCategory tests the GetRulesByCategory function with error cases
func TestGetRulesByCategory(t *testing.T) {
	tests := []struct {
		name             string
		category         string
		minExpectedRules int
	}{
		{"const category", "const", 3},
		{"func category", "func", 12},
		{"var category", "var", 17},
		{"test category", "test", 6},
		{"comment category", "comment", 7},
		{"unknown category error case", "unknown", 0},
		{"empty category error case", "", 0},
	}

	// Itération sur les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := ktn.GetRulesByCategory(tt.category)

			// Vérification des cas d'erreur
			if tt.minExpectedRules == 0 {
				// Vérification règles vides
				if len(rules) > 0 {
					t.Errorf("GetRulesByCategory(%q) returned %d rules, expected 0", tt.category, len(rules))
				}
			} else {
				// Vérification nombre minimum
				if len(rules) < tt.minExpectedRules {
					t.Errorf("GetRulesByCategory(%q) returned %d rules, expected at least %d", tt.category, len(rules), tt.minExpectedRules)
				}

				// Check that all rules are non-nil
				for i, rule := range rules {
					// Vérification règle
					if rule == nil {
						t.Errorf("Rule at index %d for category %q is nil", i, tt.category)
					}
				}
			}
		})
	}
}
