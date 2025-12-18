// Package prompt_test provides black-box tests for phase classification.
package prompt_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/prompt"
)

// TestClassifyRule_Structural tests structural rule classification.
//
// Params:
//   - t: testing object
func TestClassifyRule_Structural(t *testing.T) {
	// Test structural rules
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
}

// TestClassifyRule_TestOrg tests test organization rule classification.
//
// Params:
//   - t: testing object
func TestClassifyRule_TestOrg(t *testing.T) {
	// Test organization rules
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
}

// TestClassifyRule_Comment tests comment rule classification.
//
// Params:
//   - t: testing object
func TestClassifyRule_Comment(t *testing.T) {
	// Comment rules
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
}

// TestClassifyRule_Local tests local fix rule classification.
//
// Params:
//   - t: testing object
func TestClassifyRule_Local(t *testing.T) {
	// Local fix rules (default)
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
}

// TestGetPhaseInfo tests phase metadata retrieval.
//
// Params:
//   - t: testing object
func TestGetPhaseInfo(t *testing.T) {
	// Test all phases
	tests := []struct {
		name         string
		phase        prompt.RulePhase
		wantName     string
		wantRerun    bool
		wantDescLen  bool
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

// TestSortRulesByPhase tests phase grouping and ordering.
//
// Params:
//   - t: testing object
func TestSortRulesByPhase(t *testing.T) {
	// Create mixed rules
	rules := []prompt.RuleViolations{
		{Code: "KTN-COMMENT-001", Violations: []prompt.Violation{{}}},
		{Code: "KTN-FUNC-001", Violations: []prompt.Violation{{}}},
		{Code: "KTN-STRUCT-004", Violations: []prompt.Violation{{}}},
		{Code: "KTN-TEST-001", Violations: []prompt.Violation{{}}},
	}

	// Sort rules
	groups := prompt.SortRulesByPhase(rules)

	// Verify we have all 4 phases
	if len(groups) != 4 {
		t.Errorf("SortRulesByPhase() returned %d groups, want 4", len(groups))
		return
	}

	// Verify phase order
	expectedOrder := []prompt.RulePhase{
		prompt.PhaseStructural,
		prompt.PhaseTestOrg,
		prompt.PhaseLocal,
		prompt.PhaseComment,
	}

	// Check each phase
	for i, expected := range expectedOrder {
		if groups[i].Phase != expected {
			t.Errorf("groups[%d].Phase = %v, want %v", i, groups[i].Phase, expected)
		}
	}
}

// TestSortRulesByPhase_EmptyPhases tests that empty phases are excluded.
//
// Params:
//   - t: testing object
func TestSortRulesByPhase_EmptyPhases(t *testing.T) {
	// Create rules for only two phases
	rules := []prompt.RuleViolations{
		{Code: "KTN-FUNC-001", Violations: []prompt.Violation{{}}},
		{Code: "KTN-VAR-002", Violations: []prompt.Violation{{}}},
	}

	// Sort rules
	groups := prompt.SortRulesByPhase(rules)

	// Verify only local phase exists
	if len(groups) != 1 {
		t.Errorf("SortRulesByPhase() returned %d groups, want 1", len(groups))
		return
	}

	// Verify it's the local phase
	if groups[0].Phase != prompt.PhaseLocal {
		t.Errorf("groups[0].Phase = %v, want PhaseLocal", groups[0].Phase)
	}

	// Verify both rules are in the group
	if len(groups[0].Rules) != 2 {
		t.Errorf("groups[0].Rules has %d rules, want 2", len(groups[0].Rules))
	}
}

// TestSortRulesByPhase_SortsRulesWithinPhase tests alphabetical sorting.
//
// Params:
//   - t: testing object
func TestSortRulesByPhase_SortsRulesWithinPhase(t *testing.T) {
	// Create unsorted rules in same phase
	rules := []prompt.RuleViolations{
		{Code: "KTN-FUNC-003", Violations: []prompt.Violation{{}}},
		{Code: "KTN-FUNC-001", Violations: []prompt.Violation{{}}},
		{Code: "KTN-FUNC-002", Violations: []prompt.Violation{{}}},
	}

	// Sort rules
	groups := prompt.SortRulesByPhase(rules)

	// Verify sorting within phase
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}

	// Check order
	expected := []string{"KTN-FUNC-001", "KTN-FUNC-002", "KTN-FUNC-003"}
	for i, code := range expected {
		if groups[0].Rules[i].Code != code {
			t.Errorf("groups[0].Rules[%d].Code = %q, want %q", i, groups[0].Rules[i].Code, code)
		}
	}
}
