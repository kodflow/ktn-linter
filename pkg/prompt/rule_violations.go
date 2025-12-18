// Package prompt provides AI-optimized prompt generation for KTN linter violations.
package prompt

// RuleViolations groups all violations for a single rule.
// Contains rule metadata (code, category, description) and all violation instances.
type RuleViolations struct {
	// Code is the rule identifier (e.g., KTN-FUNC-001).
	Code string
	// Category is the rule category (e.g., func, var, struct).
	Category string
	// Description is the rule description.
	Description string
	// GoodExample is the content from good.go testdata.
	GoodExample string
	// Phase indicates the structural impact phase.
	Phase RulePhase
	// Violations contains all violations for this rule.
	Violations []Violation
}
