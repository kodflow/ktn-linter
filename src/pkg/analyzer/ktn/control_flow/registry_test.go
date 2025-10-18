package ktn_control_flow

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
			continue
		}
		if rule.Name == "" {
			t.Errorf("AllRules[%d] has empty name", i)
		}
	}
}
