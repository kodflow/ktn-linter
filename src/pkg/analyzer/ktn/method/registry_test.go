package ktn_method

import (
	"testing"
)

// TestGetRules tests the functionality of the corresponding implementation.
func TestGetRules(t *testing.T) {
	rules := GetRules()
	if rules == nil {
		t.Fatal("GetRules() should not return nil")
	}
	if len(rules) != len(AllRules) {
		t.Errorf("GetRules() returned %d rules, expected %d", len(rules), len(AllRules))
	}
}

// TestAllRulesNotNil tests the functionality of the corresponding implementation.
func TestAllRulesNotNil(t *testing.T) {
	if AllRules == nil {
		t.Fatal("AllRules should not be nil")
	}
	for i, rule := range AllRules {
		if rule == nil {
			t.Errorf("AllRules[%d] is nil", i)
		}
		if rule.Name == "" {
			t.Errorf("AllRules[%d] has empty name", i)
		}
	}
}
