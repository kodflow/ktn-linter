package ktn_test

import (
	"testing"
)

// TestGetRules vérifie que GetRules retourne toutes les règles définies.
func TestGetRules(t *testing.T) {
	rules := GetRules()
	if rules == nil {
		t.Fatal("GetRules() should not return nil")
	}
	if len(rules) != len(AllRules) {
		t.Errorf("GetRules() returned %d rules, expected %d", len(rules), len(AllRules))
	}
}

// TestAllRulesNotNil vérifie qu'aucune règle dans AllRules n'est nil.
func TestAllRulesNotNil(t *testing.T) {
	if AllRules == nil {
		t.Fatal("AllRules should not be nil")
	}
	for i, rule := range AllRules {
		if rule == nil {
			t.Errorf("AllRules[%d] is nil", i)
		}
	}
}

// TestAllRulesContainValidAnalyzers vérifie que toutes les règles ont des noms valides et uniques.
func TestAllRulesContainValidAnalyzers(t *testing.T) {
	namesSeen := make(map[string]bool)
	for _, analyzer := range AllRules {
		if analyzer.Name == "" {
			t.Error("Found analyzer with empty name")
		}
		if namesSeen[analyzer.Name] {
			t.Errorf("Duplicate analyzer name found: %s", analyzer.Name)
		}
		namesSeen[analyzer.Name] = true
	}
}
