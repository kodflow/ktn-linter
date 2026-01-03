// Package rules provides rule information extraction and formatting utilities.
package rules

// RulesOutput is the complete output structure for the rules command.
// It aggregates all rules with metadata for display purposes.
type RulesOutput struct {
	TotalCount int        // Total number of rules
	Categories []string   // Available categories
	Rules      []RuleInfo // All rules (filtered if requested)
}
