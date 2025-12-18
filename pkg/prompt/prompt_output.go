// Package prompt provides AI-optimized prompt generation for KTN linter violations.
package prompt

// PromptOutput is the complete output structure for prompt generation.
// Contains total counts and all phase groups with their rules.
type PromptOutput struct {
	// TotalViolations is the total count of violations across all rules.
	TotalViolations int
	// TotalRules is the count of rules with at least one violation.
	TotalRules int
	// Phases contains the ordered phase groups.
	Phases []PhaseGroup
}
