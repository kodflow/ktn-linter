// Package prompt provides white-box tests for internal phases functions.
package prompt

import "testing"

// Test_groupRulesByPhase tests the groupRulesByPhase private function.
//
// Params:
//   - t: testing object
func Test_groupRulesByPhase(t *testing.T) {
	// Define test cases for groupRulesByPhase
	tests := []struct {
		name           string
		rules          []RuleViolations
		expectedGroups int
		checkPhase     bool
		expectedPhase  RulePhase
	}{
		{
			name:           "empty rules",
			rules:          []RuleViolations{},
			expectedGroups: 0,
		},
		{
			name: "groups by phase",
			rules: []RuleViolations{
				{Code: "KTN-FUNC-001", Violations: []Violation{{}}},
				{Code: "KTN-STRUCT-004", Violations: []Violation{{}}},
				{Code: "KTN-TEST-001", Violations: []Violation{{}}},
				{Code: "KTN-COMMENT-001", Violations: []Violation{{}}},
			},
			expectedGroups: 4,
		},
		{
			name: "sets Phase field",
			rules: []RuleViolations{
				{Code: "KTN-FUNC-001", Violations: []Violation{{}}},
			},
			expectedGroups: 1,
			checkPhase:     true,
			expectedPhase:  PhaseLocal,
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := groupRulesByPhase(tt.rules)

			// Verify groups count
			if len(result) != tt.expectedGroups {
				t.Errorf("groupRulesByPhase() returned %d groups, want %d", len(result), tt.expectedGroups)
			}

			// Verify Phase field if needed
			if tt.checkPhase && len(tt.rules) > 0 {
				if tt.rules[0].Phase != tt.expectedPhase {
					t.Errorf("rule.Phase = %v, want %v", tt.rules[0].Phase, tt.expectedPhase)
				}
			}
		})
	}
}

// Test_sortRulesInPhases tests the sortRulesInPhases private function.
//
// Params:
//   - t: testing object
func Test_sortRulesInPhases(t *testing.T) {
	// Define test cases for sortRulesInPhases
	tests := []struct {
		name          string
		phaseMap      map[RulePhase][]RuleViolations
		expectedEmpty bool
		checkOrder    bool
		phase         RulePhase
		expectedFirst string
	}{
		{
			name:          "handles empty map",
			phaseMap:      map[RulePhase][]RuleViolations{},
			expectedEmpty: true,
		},
		{
			name: "sorts rules alphabetically",
			phaseMap: map[RulePhase][]RuleViolations{
				PhaseLocal: {
					{Code: "KTN-FUNC-003"},
					{Code: "KTN-FUNC-001"},
					{Code: "KTN-FUNC-002"},
				},
			},
			checkOrder:    true,
			phase:         PhaseLocal,
			expectedFirst: "KTN-FUNC-001",
		},
		{
			name: "sorts each phase independently",
			phaseMap: map[RulePhase][]RuleViolations{
				PhaseLocal: {
					{Code: "KTN-VAR-002"},
					{Code: "KTN-FUNC-001"},
				},
				PhaseComment: {
					{Code: "KTN-COMMENT-002"},
					{Code: "KTN-COMMENT-001"},
				},
			},
			checkOrder:    true,
			phase:         PhaseLocal,
			expectedFirst: "KTN-FUNC-001",
		},
	}

	// Run all test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortRulesInPhases(tt.phaseMap)

			// Check empty
			if tt.expectedEmpty && len(tt.phaseMap) != 0 {
				t.Error("phaseMap should remain empty")
			}

			// Check order
			if tt.checkOrder && len(tt.phaseMap[tt.phase]) > 0 {
				if tt.phaseMap[tt.phase][0].Code != tt.expectedFirst {
					t.Errorf("first code = %q, want %q", tt.phaseMap[tt.phase][0].Code, tt.expectedFirst)
				}
			}
		})
	}
}

// Test_buildPhaseGroups tests the buildPhaseGroups private function.
//
// Params:
//   - t: testing object
func Test_buildPhaseGroups(t *testing.T) {
	// Define test cases for buildPhaseGroups
	tests := []struct {
		name           string
		phaseMap       map[RulePhase][]RuleViolations
		expectedGroups int
		checkMetadata  bool
		expectedName   string
		checkRules     bool
		expectedRules  int
	}{
		{
			name:           "empty map returns empty slice",
			phaseMap:       map[RulePhase][]RuleViolations{},
			expectedGroups: 0,
		},
		{
			name: "skips empty phases",
			phaseMap: map[RulePhase][]RuleViolations{
				PhaseLocal:   {{Code: "KTN-FUNC-001", Violations: []Violation{{}}}},
				PhaseComment: {},
			},
			expectedGroups: 1,
		},
		{
			name: "maintains phase order",
			phaseMap: map[RulePhase][]RuleViolations{
				PhaseComment:    {{Code: "KTN-COMMENT-001", Violations: []Violation{{}}}},
				PhaseLocal:      {{Code: "KTN-FUNC-001", Violations: []Violation{{}}}},
				PhaseStructural: {{Code: "KTN-STRUCT-004", Violations: []Violation{{}}}},
				PhaseTestOrg:    {{Code: "KTN-TEST-001", Violations: []Violation{{}}}},
			},
			expectedGroups: 4,
		},
		{
			name: "populates phase metadata",
			phaseMap: map[RulePhase][]RuleViolations{
				PhaseStructural: {{Code: "KTN-STRUCT-004", Violations: []Violation{{}}}},
			},
			expectedGroups: 1,
			checkMetadata:  true,
			expectedName:   "Structural Changes",
		},
		{
			name: "includes rules in group",
			phaseMap: map[RulePhase][]RuleViolations{
				PhaseLocal: {
					{Code: "KTN-FUNC-001", Violations: []Violation{{}}},
					{Code: "KTN-FUNC-002", Violations: []Violation{{}}},
				},
			},
			expectedGroups: 1,
			checkRules:     true,
			expectedRules:  2,
		},
	}

	// Run all test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildPhaseGroups(tt.phaseMap)

			// Verify groups count
			if len(result) != tt.expectedGroups {
				t.Errorf("buildPhaseGroups() returned %d groups, want %d", len(result), tt.expectedGroups)
				return
			}

			// Check metadata
			if tt.checkMetadata && len(result) > 0 {
				if result[0].Name != tt.expectedName {
					t.Errorf("Name = %q, want %q", result[0].Name, tt.expectedName)
				}
			}

			// Check rules
			if tt.checkRules && len(result) > 0 {
				if len(result[0].Rules) != tt.expectedRules {
					t.Errorf("Rules count = %d, want %d", len(result[0].Rules), tt.expectedRules)
				}
			}
		})
	}
}
