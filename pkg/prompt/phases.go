// Package prompt provides AI-optimized prompt generation for KTN linter violations.
package prompt

import (
	"sort"
	"strings"
)

const (
	// phaseCount is the number of defined phases for map preallocation.
	phaseCount int = 4
)

var (
	// structuralRules defines rules that can move/delete files (must run first).
	structuralRules map[string]bool = map[string]bool{
		"KTN-STRUCT-004": true, // One struct per file -> splits files
	}

	// testOrgRules defines test file organization rules.
	testOrgRules map[string]bool = map[string]bool{
		"KTN-TEST-001": true, // Test file naming suffix
		"KTN-TEST-003": true, // Test file pairing
		"KTN-TEST-006": true, // Pattern 1:1 files test/source
		"KTN-TEST-008": true, // Test files required
		"KTN-TEST-009": true, // Public tests in _external_test.go
		"KTN-TEST-010": true, // Private tests in _internal_test.go
		"KTN-TEST-011": true, // Correct package declaration
	}

	// commentRules defines comment rules (should run last).
	commentRules map[string]bool = map[string]bool{
		"KTN-COMMENT-001": true, // Inline comment length
		"KTN-COMMENT-002": true, // Package comment
		"KTN-COMMENT-003": true, // Constant documentation
		"KTN-COMMENT-004": true, // Variable documentation
		"KTN-COMMENT-005": true, // Function documentation
		"KTN-COMMENT-006": true, // Control structure comments
		"KTN-COMMENT-007": true, // Logic comments
	}
)

// ClassifyRule determines the phase for a rule code.
//
// Params:
//   - code: the rule code (e.g., KTN-FUNC-001)
//
// Returns:
//   - RulePhase: the appropriate phase for this rule
func ClassifyRule(code string) RulePhase {
	// Check structural rules first (highest priority)
	if structuralRules[code] {
		// Return structural phase for file-modifying rules
		return PhaseStructural
	}

	// Check test organization rules
	if testOrgRules[code] {
		// Return test org phase for test file rules
		return PhaseTestOrg
	}

	// Check comment rules (lowest priority, applied last)
	if commentRules[code] {
		// Return comment phase for documentation rules
		return PhaseComment
	}

	// Check modernize rules (applied as local fixes)
	if strings.HasPrefix(code, "KTN-MODERNIZE-") {
		// Return local phase for modernize suggestions
		return PhaseLocal
	}

	// Return local phase as default
	return PhaseLocal
}

// GetPhaseInfo returns metadata for a phase.
//
// Params:
//   - phase: the phase to get info for
//
// Returns:
//   - name: human-readable phase name
//   - description: instructions for this phase
//   - needsRerun: whether linter should be re-run after
func GetPhaseInfo(phase RulePhase) (name, description string, needsRerun bool) {
	// Determine phase metadata based on phase type
	switch phase {
	// Handle structural phase
	case PhaseStructural:
		// Return info for file-modifying rules
		return "Structural Changes",
			"These rules may create, move, or delete files. Re-run linter after this phase.",
			true
	// Handle test organization phase
	case PhaseTestOrg:
		// Return info for test file rules
		return "Test Organization",
			"Test file naming and placement rules. May require file renames.",
			true
	// Handle local fixes phase
	case PhaseLocal:
		// Return info for in-file modifications
		return "Local Fixes",
			"Code modifications within existing files. No structural changes.",
			false
	// Handle comments phase
	case PhaseComment:
		// Return info for documentation rules
		return "Comments & Documentation",
			"Documentation rules. Apply last after all code is finalized.",
			false
	// Handle unknown phase
	default:
		// Return empty info for unknown phase
		return "Unknown", "", false
	}
}

// SortRulesByPhase organizes rules into their execution phases.
//
// Params:
//   - rules: slice of rule violations to organize
//
// Returns:
//   - []PhaseGroup: ordered phase groups with rules
func SortRulesByPhase(rules []RuleViolations) []PhaseGroup {
	// Group rules by phase
	phaseMap := groupRulesByPhase(rules)

	// Sort rules within each phase
	sortRulesInPhases(phaseMap)

	// Build ordered phase groups
	return buildPhaseGroups(phaseMap)
}

// groupRulesByPhase groups rules by their classified phase.
// Groups each rule into the appropriate execution phase based on its code.
//
// Params:
//   - rules: slice of rule violations
//
// Returns:
//   - map[RulePhase][]RuleViolations: rules grouped by phase
func groupRulesByPhase(rules []RuleViolations) map[RulePhase][]RuleViolations {
	// Preallocate map with known phase count
	phaseMap := make(map[RulePhase][]RuleViolations, phaseCount)

	// Classify and group each rule
	for i := range rules {
		rules[i].Phase = ClassifyRule(rules[i].Code)
		phaseMap[rules[i].Phase] = append(phaseMap[rules[i].Phase], rules[i])
	}

	// Return grouped rules
	return phaseMap
}

// sortRulesInPhases sorts rules within each phase by code.
//
// Params:
//   - phaseMap: map of phases to rules
func sortRulesInPhases(phaseMap map[RulePhase][]RuleViolations) {
	// Sort rules in each phase by code
	for phase := range phaseMap {
		sort.Slice(phaseMap[phase], func(i, j int) bool {
			// Compare rule codes alphabetically
			return phaseMap[phase][i].Code < phaseMap[phase][j].Code
		})
	}
}

// buildPhaseGroups constructs ordered phase groups from the phase map.
//
// Params:
//   - phaseMap: map of phases to rules
//
// Returns:
//   - []PhaseGroup: ordered phase groups
func buildPhaseGroups(phaseMap map[RulePhase][]RuleViolations) []PhaseGroup {
	// Define phase order
	phases := []RulePhase{PhaseStructural, PhaseTestOrg, PhaseLocal, PhaseComment}
	var result []PhaseGroup

	// Build group for each phase with rules
	for _, phase := range phases {
		rules, exists := phaseMap[phase]
		// Skip empty phases
		if !exists || len(rules) == 0 {
			continue
		}

		// Get phase metadata
		name, desc, needsRerun := GetPhaseInfo(phase)
		result = append(result, PhaseGroup{
			Phase:       phase,
			Name:        name,
			Description: desc,
			Rules:       rules,
			NeedsRerun:  needsRerun,
		})
	}

	// Return ordered phase groups
	return result
}
