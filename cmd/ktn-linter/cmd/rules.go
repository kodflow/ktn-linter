// Package cmd implements the CLI commands for ktn-linter.
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/rules"
	"github.com/spf13/cobra"
)

const (
	// flagRulesFormat is the flag name for output format.
	flagRulesFormat string = "format"
	// flagRulesNoExamples is the flag name to skip examples.
	flagRulesNoExamples string = "no-examples"
	// defaultRulesFormat is the default output format.
	defaultRulesFormat string = "markdown"
)

// rulesCmd represents the rules command.
var rulesCmd *cobra.Command = &cobra.Command{
	Use:   "rules",
	Short: "Display all available KTN rules",
	Long:  `Display all KTN lint rules with descriptions and code examples for AI-assisted development.`,
	Run:   runRules,
}

// init registers the rules command with root.
//
// Params: none
//
// Returns: none
func init() {
	rootCmd.AddCommand(rulesCmd)
	// Add format flag
	rulesCmd.Flags().String(flagRulesFormat, defaultRulesFormat, "Output format: markdown, json, text")
	// Add no-examples flag
	rulesCmd.Flags().Bool(flagRulesNoExamples, false, "Exclude code examples from output")
}

// runRules executes the rules command.
//
// Params:
//   - cmd: Cobra command (used to get flags)
//   - args: command line arguments (unused)
//
// Returns: none
func runRules(cmd *cobra.Command, _ []string) {
	// Parse command-specific flags
	opts := parseRulesOptions(cmd)

	// Get rules based on filters from root flags
	infos := getRulesWithFilters()

	// Load examples if requested
	if !opts.NoExamples {
		infos = rules.LoadGoodExamples(infos)
	}

	// Build output structure
	output := buildRulesOutput(infos)

	// Format and display
	formatRulesOutput(output, opts.Format)
}

// rulesOptions contains options for the rules command.
type rulesOptions struct {
	Format     string // Output format (markdown, json, text)
	NoExamples bool   // Whether to skip examples
}

// parseRulesOptions extracts options from command flags.
//
// Params:
//   - cmd: Cobra command with flags
//
// Returns:
//   - rulesOptions: extracted options
func parseRulesOptions(cmd *cobra.Command) rulesOptions {
	format, _ := cmd.Flags().GetString(flagRulesFormat)
	noExamples, _ := cmd.Flags().GetBool(flagRulesNoExamples)

	// Return parsed options
	return rulesOptions{
		Format:     format,
		NoExamples: noExamples,
	}
}

// getRulesWithFilters returns rules applying category/code filters from root.
//
// Returns:
//   - []rules.RuleInfo: filtered rules
func getRulesWithFilters() []rules.RuleInfo {
	flags := rootCmd.PersistentFlags()
	category, _ := flags.GetString(flagCategory)
	onlyRule, _ := flags.GetString(flagOnlyRule)

	// Filter by specific rule code
	if onlyRule != "" {
		info := rules.GetRuleInfoByCode(onlyRule)
		// Check if found
		if info == nil {
			fmt.Fprintf(os.Stderr, "Rule not found: %s\n", onlyRule)
			OsExit(1)
		}
		// Return single rule in slice
		return []rules.RuleInfo{*info}
	}

	// Filter by category
	if category != "" {
		// Return rules for specified category
		return rules.GetRuleInfosByCategory(category)
	}

	// Return all available rules
	return rules.GetAllRuleInfos()
}

// buildRulesOutput creates the output structure from rule infos.
//
// Params:
//   - infos: rule information slice
//
// Returns:
//   - rules.RulesOutput: structured output
func buildRulesOutput(infos []rules.RuleInfo) rules.RulesOutput {
	// Return output structure
	return rules.RulesOutput{
		TotalCount: len(infos),
		Categories: rules.GetCategories(),
		Rules:      infos,
	}
}

// formatRulesOutput formats and displays the rules output.
//
// Params:
//   - output: structured rules output
//   - format: output format (markdown, json, text)
//
// Returns: none
func formatRulesOutput(output rules.RulesOutput, format string) {
	// Select formatter based on format
	switch strings.ToLower(format) {
	// JSON format output
	case "json":
		formatRulesJSON(os.Stdout, output)
	// Plain text format output
	case "text":
		formatRulesText(os.Stdout, output)
	// Default to markdown format
	default:
		formatRulesMarkdown(os.Stdout, output)
	}
}

// formatRulesMarkdown outputs rules in Markdown format.
//
// Params:
//   - w: writer to output to
//   - output: rules output structure
//
// Returns: none
func formatRulesMarkdown(w io.Writer, output rules.RulesOutput) {
	// Write header
	fmt.Fprintf(w, "# KTN-Linter Rules Reference\n\n")
	fmt.Fprintf(w, "**Total**: %d rules | **Categories**: %s\n\n",
		output.TotalCount, strings.Join(output.Categories, ", "))
	fmt.Fprintf(w, "---\n\n")

	// Write each rule
	for _, rule := range output.Rules {
		writeRuleMarkdown(w, rule)
	}
}

// writeRuleMarkdown writes a single rule in Markdown format.
//
// Params:
//   - w: writer to output to
//   - rule: rule information
//
// Returns: none
func writeRuleMarkdown(w io.Writer, rule rules.RuleInfo) {
	// Write rule header
	fmt.Fprintf(w, "## %s\n\n", rule.Code)
	fmt.Fprintf(w, "**Category**: %s\n\n", rule.Category)
	fmt.Fprintf(w, "%s\n\n", rule.Description)

	// Write example if available
	if rule.GoodExample != "" {
		fmt.Fprintf(w, "### Good Example\n\n")
		fmt.Fprintf(w, "```go\n%s```\n\n", rule.GoodExample)
	}

	// Write separator
	fmt.Fprintf(w, "---\n\n")
}

// formatRulesJSON outputs rules in JSON format.
//
// Params:
//   - w: writer to output to
//   - output: rules output structure
//
// Returns: none
func formatRulesJSON(w io.Writer, output rules.RulesOutput) {
	// Create JSON encoder with indentation
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	// Encode output
	if err := encoder.Encode(output); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		OsExit(1)
	}
}

// formatRulesText outputs rules in plain text format.
//
// Params:
//   - w: writer to output to
//   - output: rules output structure
//
// Returns: none
func formatRulesText(w io.Writer, output rules.RulesOutput) {
	// Write header
	fmt.Fprintf(w, "KTN Rules - %d rules in %d categories\n\n",
		output.TotalCount, len(output.Categories))

	// Group by category
	currentCategory := ""
	// Write each rule
	for _, rule := range output.Rules {
		// Check if category changed
		if rule.Category != currentCategory {
			currentCategory = rule.Category
			fmt.Fprintf(w, "=== %s ===\n\n", strings.ToUpper(currentCategory))
		}
		// Write rule
		writeRuleText(w, rule)
	}
}

// writeRuleText writes a single rule in plain text format.
//
// Params:
//   - w: writer to output to
//   - rule: rule information
//
// Returns: none
func writeRuleText(w io.Writer, rule rules.RuleInfo) {
	// Write rule code and description
	fmt.Fprintf(w, "%s\n", rule.Code)
	fmt.Fprintf(w, "  %s\n", rule.Description)

	// Write example if available
	if rule.GoodExample != "" {
		fmt.Fprintf(w, "\n  Good Example:\n")
		// Iterate using SplitSeq for efficiency
		for line := range strings.SplitSeq(rule.GoodExample, "\n") {
			fmt.Fprintf(w, "    %s\n", line)
		}
	}

	// Write separator
	fmt.Fprintf(w, "\n---\n\n")
}
