// Package prompt_test provides black-box tests for phase classification.
package prompt_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/prompt"
)

// TestClassifyRule tests the ClassifyRule public function.
//
// Params:
//   - t: testing object
func TestClassifyRule(t *testing.T) {
	// Test structural rules
	t.Run("structural", func(t *testing.T) {
		tests := []struct {
			name string
			code string
			want prompt.RulePhase
		}{
			{
				name: "struct-004 is structural",
				code: "KTN-STRUCT-004",
				want: prompt.PhaseStructural,
			},
		}

		// Run test cases
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := prompt.ClassifyRule(tt.code)
				// Verify classification
				if got != tt.want {
					t.Errorf("ClassifyRule(%q) = %v, want %v", tt.code, got, tt.want)
				}
			})
		}
	})

	// Test organization rules
	t.Run("test org", func(t *testing.T) {
		tests := []struct {
			name string
			code string
			want prompt.RulePhase
		}{
			{name: "test-001", code: "KTN-TEST-001", want: prompt.PhaseTestOrg},
			{name: "test-003", code: "KTN-TEST-003", want: prompt.PhaseTestOrg},
			{name: "test-006", code: "KTN-TEST-006", want: prompt.PhaseTestOrg},
			{name: "test-008", code: "KTN-TEST-008", want: prompt.PhaseTestOrg},
			{name: "test-009", code: "KTN-TEST-009", want: prompt.PhaseTestOrg},
			{name: "test-010", code: "KTN-TEST-010", want: prompt.PhaseTestOrg},
			{name: "test-011", code: "KTN-TEST-011", want: prompt.PhaseTestOrg},
		}

		// Run test cases
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := prompt.ClassifyRule(tt.code)
				// Verify classification
				if got != tt.want {
					t.Errorf("ClassifyRule(%q) = %v, want %v", tt.code, got, tt.want)
				}
			})
		}
	})

	// Comment rules
	t.Run("comment", func(t *testing.T) {
		tests := []struct {
			name string
			code string
			want prompt.RulePhase
		}{
			{name: "comment-001", code: "KTN-COMMENT-001", want: prompt.PhaseComment},
			{name: "comment-002", code: "KTN-COMMENT-002", want: prompt.PhaseComment},
			{name: "comment-003", code: "KTN-COMMENT-003", want: prompt.PhaseComment},
			{name: "comment-004", code: "KTN-COMMENT-004", want: prompt.PhaseComment},
			{name: "comment-005", code: "KTN-COMMENT-005", want: prompt.PhaseComment},
			{name: "comment-006", code: "KTN-COMMENT-006", want: prompt.PhaseComment},
			{name: "comment-007", code: "KTN-COMMENT-007", want: prompt.PhaseComment},
		}

		// Run test cases
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := prompt.ClassifyRule(tt.code)
				// Verify classification
				if got != tt.want {
					t.Errorf("ClassifyRule(%q) = %v, want %v", tt.code, got, tt.want)
				}
			})
		}
	})

	// Local fix rules (default)
	t.Run("local", func(t *testing.T) {
		tests := []struct {
			name string
			code string
			want prompt.RulePhase
		}{
			{name: "func-001", code: "KTN-FUNC-001", want: prompt.PhaseLocal},
			{name: "var-002", code: "KTN-VAR-002", want: prompt.PhaseLocal},
			{name: "struct-001", code: "KTN-STRUCT-001", want: prompt.PhaseLocal},
			{name: "const-001", code: "KTN-CONST-001", want: prompt.PhaseLocal},
			{name: "return-001", code: "KTN-RETURN-001", want: prompt.PhaseLocal},
			{name: "interface-001", code: "KTN-INTERFACE-001", want: prompt.PhaseLocal},
		}

		// Run test cases
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := prompt.ClassifyRule(tt.code)
				// Verify classification
				if got != tt.want {
					t.Errorf("ClassifyRule(%q) = %v, want %v", tt.code, got, tt.want)
				}
			})
		}
	})
}

// TestGetPhaseInfo tests phase metadata retrieval.
//
// Params:
//   - t: testing object
func TestGetPhaseInfo(t *testing.T) {
	// Test all phases
	tests := []struct {
		name        string
		phase       prompt.RulePhase
		wantName    string
		wantRerun   bool
		wantDescLen bool
	}{
		{
			name:        "structural",
			phase:       prompt.PhaseStructural,
			wantName:    "Structural Changes",
			wantRerun:   true,
			wantDescLen: true,
		},
		{
			name:        "test org",
			phase:       prompt.PhaseTestOrg,
			wantName:    "Test Organization",
			wantRerun:   true,
			wantDescLen: true,
		},
		{
			name:        "local",
			phase:       prompt.PhaseLocal,
			wantName:    "Local Fixes",
			wantRerun:   false,
			wantDescLen: true,
		},
		{
			name:        "comment",
			phase:       prompt.PhaseComment,
			wantName:    "Comments & Documentation",
			wantRerun:   false,
			wantDescLen: true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, desc, needsRerun := prompt.GetPhaseInfo(tt.phase)

			// Verify name
			if name != tt.wantName {
				t.Errorf("GetPhaseInfo(%v) name = %q, want %q", tt.phase, name, tt.wantName)
			}

			// Verify needs rerun
			if needsRerun != tt.wantRerun {
				t.Errorf("GetPhaseInfo(%v) needsRerun = %v, want %v", tt.phase, needsRerun, tt.wantRerun)
			}

			// Verify description exists
			if tt.wantDescLen && len(desc) == 0 {
				t.Errorf("GetPhaseInfo(%v) description is empty", tt.phase)
			}
		})
	}
}

// TestGetPhaseInfo_UnknownPhase tests handling of unknown phase values.
//
// Params:
//   - t: testing object
func TestGetPhaseInfo_UnknownPhase(t *testing.T) {
	// Define test cases for unknown phase handling
	tests := []struct {
		name           string
		phase          prompt.RulePhase
		expectedName   string
		expectedDesc   string
		expectedRerun  bool
	}{
		{
			name:          "unknown phase returns default values",
			phase:         prompt.RulePhase(999),
			expectedName:  "Unknown",
			expectedDesc:  "",
			expectedRerun: false,
		},
	}

	// Run all test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get phase info
			name, desc, needsRerun := prompt.GetPhaseInfo(tt.phase)

			// Verify name
			if name != tt.expectedName {
				t.Errorf("GetPhaseInfo() name = %q, want %q", name, tt.expectedName)
			}

			// Verify description
			if desc != tt.expectedDesc {
				t.Errorf("GetPhaseInfo() description = %q, want %q", desc, tt.expectedDesc)
			}

			// Verify needsRerun
			if needsRerun != tt.expectedRerun {
				t.Errorf("GetPhaseInfo() needsRerun = %v, want %v", needsRerun, tt.expectedRerun)
			}
		})
	}
}

// TestSortRulesByPhase tests phase grouping and ordering.
//
// Params:
//   - t: testing object
func TestSortRulesByPhase(t *testing.T) {
	// Define test cases for phase sorting
	tests := []struct {
		name           string
		rules          []prompt.RuleViolations
		expectedPhases []prompt.RulePhase
	}{
		{
			name: "groups and orders mixed rules correctly",
			rules: []prompt.RuleViolations{
				{Code: "KTN-COMMENT-001", Violations: []prompt.Violation{{}}},
				{Code: "KTN-FUNC-001", Violations: []prompt.Violation{{}}},
				{Code: "KTN-STRUCT-004", Violations: []prompt.Violation{{}}},
				{Code: "KTN-TEST-001", Violations: []prompt.Violation{{}}},
			},
			expectedPhases: []prompt.RulePhase{
				prompt.PhaseStructural,
				prompt.PhaseTestOrg,
				prompt.PhaseLocal,
				prompt.PhaseComment,
			},
		},
	}

	// Run all test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Sort rules
			groups := prompt.SortRulesByPhase(tt.rules)

			// Verify group count
			if len(groups) != len(tt.expectedPhases) {
				t.Errorf("SortRulesByPhase() returned %d groups, want %d", len(groups), len(tt.expectedPhases))
				return
			}

			// Check each phase
			for i, expected := range tt.expectedPhases {
				if groups[i].Phase != expected {
					t.Errorf("groups[%d].Phase = %v, want %v", i, groups[i].Phase, expected)
				}
			}
		})
	}
}

// TestSortRulesByPhase_EmptyPhases tests that empty phases are excluded.
//
// Params:
//   - t: testing object
func TestSortRulesByPhase_EmptyPhases(t *testing.T) {
	// Define test cases for empty phase handling
	tests := []struct {
		name           string
		rules          []prompt.RuleViolations
		expectedGroups int
		expectedPhase  prompt.RulePhase
		expectedRules  int
	}{
		{
			name: "excludes phases without rules",
			rules: []prompt.RuleViolations{
				{Code: "KTN-FUNC-001", Violations: []prompt.Violation{{}}},
				{Code: "KTN-VAR-002", Violations: []prompt.Violation{{}}},
			},
			expectedGroups: 1,
			expectedPhase:  prompt.PhaseLocal,
			expectedRules:  2,
		},
	}

	// Run all test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Sort rules
			groups := prompt.SortRulesByPhase(tt.rules)

			// Verify group count
			if len(groups) != tt.expectedGroups {
				t.Errorf("SortRulesByPhase() returned %d groups, want %d", len(groups), tt.expectedGroups)
				return
			}

			// Verify phase
			if groups[0].Phase != tt.expectedPhase {
				t.Errorf("groups[0].Phase = %v, want %v", groups[0].Phase, tt.expectedPhase)
			}

			// Verify rule count
			if len(groups[0].Rules) != tt.expectedRules {
				t.Errorf("groups[0].Rules has %d rules, want %d", len(groups[0].Rules), tt.expectedRules)
			}
		})
	}
}

// TestSortRulesByPhase_SortsRulesWithinPhase tests alphabetical sorting.
//
// Params:
//   - t: testing object
func TestSortRulesByPhase_SortsRulesWithinPhase(t *testing.T) {
	// Define test cases for rule sorting within phase
	tests := []struct {
		name          string
		rules         []prompt.RuleViolations
		expectedCodes []string
	}{
		{
			name: "sorts rules alphabetically within phase",
			rules: []prompt.RuleViolations{
				{Code: "KTN-FUNC-003", Violations: []prompt.Violation{{}}},
				{Code: "KTN-FUNC-001", Violations: []prompt.Violation{{}}},
				{Code: "KTN-FUNC-002", Violations: []prompt.Violation{{}}},
			},
			expectedCodes: []string{"KTN-FUNC-001", "KTN-FUNC-002", "KTN-FUNC-003"},
		},
	}

	// Run all test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Sort rules
			groups := prompt.SortRulesByPhase(tt.rules)

			// Verify single group
			if len(groups) != 1 {
				t.Fatalf("expected 1 group, got %d", len(groups))
			}

			// Check order
			for i, code := range tt.expectedCodes {
				if groups[0].Rules[i].Code != code {
					t.Errorf("groups[0].Rules[%d].Code = %q, want %q", i, groups[0].Rules[i].Code, code)
				}
			}
		})
	}
}
