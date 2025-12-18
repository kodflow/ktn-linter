// Package prompt provides AI-optimized prompt generation for KTN linter violations.
package prompt

// RulePhase categorizes rules by their structural impact on the codebase.
type RulePhase int

const (
	// PhaseStructural contains rules that can move/delete files.
	PhaseStructural RulePhase = iota
	// PhaseTestOrg contains test file organization rules.
	PhaseTestOrg
	// PhaseLocal contains code modification rules.
	PhaseLocal
	// PhaseComment contains comment/documentation rules (applied last).
	PhaseComment
)
