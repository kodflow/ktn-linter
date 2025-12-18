// Package prompt provides AI-optimized prompt generation for KTN linter violations.
package prompt

// PhaseGroup groups rules by their execution phase.
// Contains phase metadata and all rules belonging to this phase.
type PhaseGroup struct {
	// Phase is the phase identifier.
	Phase RulePhase
	// Name is the human-readable phase name.
	Name string
	// Description provides instructions for this phase.
	Description string
	// Rules contains all rules in this phase.
	Rules []RuleViolations
	// NeedsRerun indicates if linter should be re-run after this phase.
	NeedsRerun bool
}
