// Package prompt provides white-box tests for internal phases functions.
package prompt

import "testing"

// Test_groupRulesByPhase tests the groupRulesByPhase private function.
//
// Params:
//   - t: testing object
func Test_groupRulesByPhase(t *testing.T) {
	// Test with empty rules
	t.Run("empty rules", func(t *testing.T) {
		result := groupRulesByPhase([]RuleViolations{})

		// Verify empty result
		if len(result) != 0 {
			t.Errorf("groupRulesByPhase([]) returned %d groups, want 0", len(result))
		}
	})

	// Test with mixed phases
	t.Run("groups by phase", func(t *testing.T) {
		rules := []RuleViolations{
			{Code: "KTN-FUNC-001", Violations: []Violation{{}}},
			{Code: "KTN-STRUCT-004", Violations: []Violation{{}}},
			{Code: "KTN-TEST-001", Violations: []Violation{{}}},
			{Code: "KTN-COMMENT-001", Violations: []Violation{{}}},
		}

		result := groupRulesByPhase(rules)

		// Verify 4 phase groups
		if len(result) != 4 {
			t.Errorf("groupRulesByPhase() returned %d groups, want 4", len(result))
		}

		// Verify each phase has rules
		if len(result[PhaseLocal]) != 1 {
			t.Errorf("PhaseLocal has %d rules, want 1", len(result[PhaseLocal]))
		}
		// Verify structural phase
		if len(result[PhaseStructural]) != 1 {
			t.Errorf("PhaseStructural has %d rules, want 1", len(result[PhaseStructural]))
		}
		// Verify test org phase
		if len(result[PhaseTestOrg]) != 1 {
			t.Errorf("PhaseTestOrg has %d rules, want 1", len(result[PhaseTestOrg]))
		}
		// Verify comment phase
		if len(result[PhaseComment]) != 1 {
			t.Errorf("PhaseComment has %d rules, want 1", len(result[PhaseComment]))
		}
	})

	// Test classifies rule phase
	t.Run("sets Phase field", func(t *testing.T) {
		rules := []RuleViolations{
			{Code: "KTN-FUNC-001", Violations: []Violation{{}}},
		}

		groupRulesByPhase(rules)

		// Verify Phase was set
		if rules[0].Phase != PhaseLocal {
			t.Errorf("rule.Phase = %v, want PhaseLocal", rules[0].Phase)
		}
	})
}

// Test_sortRulesInPhases tests the sortRulesInPhases private function.
//
// Params:
//   - t: testing object
func Test_sortRulesInPhases(t *testing.T) {
	// Test sorting within phase
	t.Run("sorts rules alphabetically", func(t *testing.T) {
		phaseMap := map[RulePhase][]RuleViolations{
			PhaseLocal: {
				{Code: "KTN-FUNC-003"},
				{Code: "KTN-FUNC-001"},
				{Code: "KTN-FUNC-002"},
			},
		}

		sortRulesInPhases(phaseMap)

		// Verify order
		expected := []string{"KTN-FUNC-001", "KTN-FUNC-002", "KTN-FUNC-003"}
		for i, code := range expected {
			// Check each position
			if phaseMap[PhaseLocal][i].Code != code {
				t.Errorf("phaseMap[PhaseLocal][%d].Code = %q, want %q",
					i, phaseMap[PhaseLocal][i].Code, code)
			}
		}
	})

	// Test empty phase map
	t.Run("handles empty map", func(t *testing.T) {
		phaseMap := map[RulePhase][]RuleViolations{}

		// Should not panic
		sortRulesInPhases(phaseMap)

		// Verify still empty
		if len(phaseMap) != 0 {
			t.Errorf("phaseMap should remain empty")
		}
	})

	// Test multiple phases
	t.Run("sorts each phase independently", func(t *testing.T) {
		phaseMap := map[RulePhase][]RuleViolations{
			PhaseLocal: {
				{Code: "KTN-VAR-002"},
				{Code: "KTN-FUNC-001"},
			},
			PhaseComment: {
				{Code: "KTN-COMMENT-002"},
				{Code: "KTN-COMMENT-001"},
			},
		}

		sortRulesInPhases(phaseMap)

		// Verify local phase order
		if phaseMap[PhaseLocal][0].Code != "KTN-FUNC-001" {
			t.Errorf("PhaseLocal[0] = %q, want KTN-FUNC-001", phaseMap[PhaseLocal][0].Code)
		}
		// Verify comment phase order
		if phaseMap[PhaseComment][0].Code != "KTN-COMMENT-001" {
			t.Errorf("PhaseComment[0] = %q, want KTN-COMMENT-001", phaseMap[PhaseComment][0].Code)
		}
	})
}

// Test_buildPhaseGroups tests the buildPhaseGroups private function.
//
// Params:
//   - t: testing object
func Test_buildPhaseGroups(t *testing.T) {
	// Test empty map
	t.Run("empty map returns empty slice", func(t *testing.T) {
		result := buildPhaseGroups(map[RulePhase][]RuleViolations{})

		// Verify empty result
		if len(result) != 0 {
			t.Errorf("buildPhaseGroups({}) returned %d groups, want 0", len(result))
		}
	})

	// Test skips empty phases
	t.Run("skips empty phases", func(t *testing.T) {
		phaseMap := map[RulePhase][]RuleViolations{
			PhaseLocal:   {{Code: "KTN-FUNC-001", Violations: []Violation{{}}}},
			PhaseComment: {}, // Empty phase
		}

		result := buildPhaseGroups(phaseMap)

		// Verify only non-empty phase included
		if len(result) != 1 {
			t.Errorf("buildPhaseGroups() returned %d groups, want 1", len(result))
		}
	})

	// Test phase order
	t.Run("maintains phase order", func(t *testing.T) {
		phaseMap := map[RulePhase][]RuleViolations{
			PhaseComment:    {{Code: "KTN-COMMENT-001", Violations: []Violation{{}}}},
			PhaseLocal:      {{Code: "KTN-FUNC-001", Violations: []Violation{{}}}},
			PhaseStructural: {{Code: "KTN-STRUCT-004", Violations: []Violation{{}}}},
			PhaseTestOrg:    {{Code: "KTN-TEST-001", Violations: []Violation{{}}}},
		}

		result := buildPhaseGroups(phaseMap)

		// Verify order
		if len(result) != 4 {
			t.Fatalf("expected 4 groups, got %d", len(result))
		}

		// Verify expected order: Structural, TestOrg, Local, Comment
		expectedOrder := []RulePhase{PhaseStructural, PhaseTestOrg, PhaseLocal, PhaseComment}
		for i, expected := range expectedOrder {
			// Check each position
			if result[i].Phase != expected {
				t.Errorf("result[%d].Phase = %v, want %v", i, result[i].Phase, expected)
			}
		}
	})

	// Test metadata population
	t.Run("populates phase metadata", func(t *testing.T) {
		phaseMap := map[RulePhase][]RuleViolations{
			PhaseStructural: {{Code: "KTN-STRUCT-004", Violations: []Violation{{}}}},
		}

		result := buildPhaseGroups(phaseMap)

		// Verify metadata
		if result[0].Name != "Structural Changes" {
			t.Errorf("Name = %q, want 'Structural Changes'", result[0].Name)
		}
		// Verify needs rerun
		if !result[0].NeedsRerun {
			t.Error("NeedsRerun should be true for structural phase")
		}
		// Verify description
		if result[0].Description == "" {
			t.Error("Description should not be empty")
		}
	})

	// Test rules are included
	t.Run("includes rules in group", func(t *testing.T) {
		phaseMap := map[RulePhase][]RuleViolations{
			PhaseLocal: {
				{Code: "KTN-FUNC-001", Violations: []Violation{{}}},
				{Code: "KTN-FUNC-002", Violations: []Violation{{}}},
			},
		}

		result := buildPhaseGroups(phaseMap)

		// Verify rules count
		if len(result[0].Rules) != 2 {
			t.Errorf("Rules count = %d, want 2", len(result[0].Rules))
		}
	})
}
