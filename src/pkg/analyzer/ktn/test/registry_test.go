package ktn_test

import (
	"testing"
)

func TestGetRules(t *testing.T) {
	rules := GetRules()
	if rules == nil {
		t.Fatal("GetRules() should not return nil")
	}
	if len(rules) != len(AllRules) {
		t.Errorf("GetRules() returned %d rules, expected %d", len(rules), len(AllRules))
	}
}

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
