// Package cmd implements the CLI commands for ktn-linter.
package cmd

import (
	"fmt"
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
	defaultRulesFormat string = "text"
	// headingSeparatorWidth is the width of heading separator lines.
	headingSeparatorWidth int = 20
	// ruleNumberDigitCount is the expected number of digits in a rule number.
	ruleNumberDigitCount int = 3
	// maxNonCategoryArgs is the maximum non-category arguments allowed.
	maxNonCategoryArgs int = 2
)

// rulesCmd represents the rules command.
var rulesCmd *cobra.Command = &cobra.Command{
	Use:   "rules [category] [rule-number]",
	Short: "Display available KTN rules",
	Long: `Display KTN lint rules in a hierarchical structure.

Usage:
  ktn-linter rules                    List all categories
  ktn-linter rules func               List rules in 'func' category
  ktn-linter rules func 001           Show details for KTN-FUNC-001
  ktn-linter rules KTN-FUNC-001       Show details for a specific rule`,
	Args: cobra.MaximumNArgs(maxNonCategoryArgs),
	Run:  runRules,
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

// RulesFormatter defines the interface for formatting rules output.
// Implementations handle different output formats (text, markdown, json).
type RulesFormatter interface {
	// DisplayCategories shows all available categories.
	DisplayCategories(categories []string)
	// DisplayCategoryRules shows all rules in a category.
	DisplayCategoryRules(category string, catRules []rules.RuleInfo)
	// DisplayRuleDetails shows detailed information for a single rule.
	DisplayRuleDetails(info rules.RuleInfo)
}

// NewRulesFormatter creates a formatter for the specified format.
//
// Params:
//   - format: output format (text, markdown, json)
//
// Returns:
//   - RulesFormatter: formatter implementation
func NewRulesFormatter(format string) RulesFormatter {
	// Select format implementation
	switch strings.ToLower(format) {
	// JSON format
	case "json":
		// Return JSON formatter
		return &jsonRulesFormatter{}
	// Markdown format
	case "markdown", "md":
		// Return Markdown formatter
		return &markdownRulesFormatter{}
	// Text format (default)
	default:
		// Return default text formatter
		return &textRulesFormatter{}
	}
}

// runRules executes the rules command (Cobra callback).
//
// Params:
//   - cmd: Cobra command (used to get flags)
//   - args: command line arguments (category and/or rule number)
//
// Returns: none
func runRules(cmd *cobra.Command, args []string) {
	// Delegate to testable function with interfaces
	runRulesWithFlags(cmd.Flags(), args)
}

// runRulesWithFlags executes the rules command with injected dependencies.
//
// Params:
//   - flags: flag getter for command flags
//   - args: command line arguments (category and/or rule number)
//
// Returns: none
func runRulesWithFlags(flags flagGetter, args []string) {
	// Parse command-specific flags
	opts := parseRulesOptions(flags)

	// Create formatter for selected format
	formatter := NewRulesFormatter(opts.Format)

	// Handle hierarchical navigation based on args
	switch len(args) {
	// No args: list categories
	case 0:
		displayCategories(formatter)
	// One arg: category name or full rule code
	case 1:
		handleSingleArg(args[0], formatter, opts)
	// Two args: category and rule number
	case maxNonCategoryArgs:
		handleCategoryAndRule(args[0], args[1], formatter, opts)
	}
}

// displayCategories shows all available categories.
//
// Params:
//   - formatter: rules formatter to use
func displayCategories(formatter RulesFormatter) {
	categories := rules.GetCategories()
	// Delegate to formatter
	formatter.DisplayCategories(categories)
}

// handleSingleArg handles one argument (category or full rule code).
//
// Params:
//   - arg: category name or rule code (e.g., "func" or "KTN-FUNC-001")
//   - formatter: rules formatter to use
//   - opts: rules options
func handleSingleArg(arg string, formatter RulesFormatter, opts rulesOptions) {
	// Check if it's a full rule code (KTN-XXX-YYY)
	if strings.HasPrefix(strings.ToUpper(arg), "KTN-") {
		// Display single rule details
		displayRuleDetails(strings.ToUpper(arg), formatter, opts)
		// Done processing full rule code
		return
	}

	// Otherwise treat as category
	displayCategoryRules(arg, formatter)
}

// displayCategoryRules shows all rules in a category.
//
// Params:
//   - category: category name
//   - formatter: rules formatter to use
func displayCategoryRules(category string, formatter RulesFormatter) {
	catRules := rules.GetRuleInfosByCategory(category)

	// Check if category exists
	if len(catRules) == 0 {
		fmt.Fprintf(os.Stderr, "Category not found: %s\n", category)
		fmt.Fprintf(os.Stderr, "Available categories: %s\n", strings.Join(rules.GetCategories(), ", "))
		OsExit(1)
	}

	// Delegate to formatter
	formatter.DisplayCategoryRules(category, catRules)
}

// handleCategoryAndRule handles two arguments (category and rule number).
//
// Params:
//   - category: category name (e.g., "func")
//   - ruleNum: rule number (e.g., "001")
//   - formatter: rules formatter to use
//   - opts: rules options
func handleCategoryAndRule(category, ruleNum string, formatter RulesFormatter, opts rulesOptions) {
	// Validate rule number format (must be 3 digits)
	if !isValidRuleNumber(ruleNum) {
		fmt.Fprintf(os.Stderr, "Invalid rule number format: %s (expected 3 digits, e.g., 001)\n", ruleNum)
		OsExit(1)
	}

	// Build full rule code
	code := fmt.Sprintf("KTN-%s-%s", strings.ToUpper(category), ruleNum)
	displayRuleDetails(code, formatter, opts)
}

// isValidRuleNumber checks if a rule number has the correct format.
// Rule numbers must be exactly 3 digits (e.g., "001", "012", "123").
//
// Params:
//   - ruleNum: rule number to validate
//
// Returns:
//   - bool: true if valid format
func isValidRuleNumber(ruleNum string) bool {
	// Check length matches expected digit count
	if len(ruleNum) != ruleNumberDigitCount {
		// Invalid length for rule number
		return false
	}

	// Check all characters are digits
	for _, c := range ruleNum {
		// Validate character is a digit
		if c < '0' || c > '9' {
			// Non-digit character found
			return false
		}
	}

	// Valid format - all checks passed
	return true
}

// displayRuleDetails shows detailed information for a single rule.
//
// Params:
//   - code: full rule code (e.g., "KTN-FUNC-001")
//   - formatter: rules formatter to use
//   - opts: rules options
func displayRuleDetails(code string, formatter RulesFormatter, opts rulesOptions) {
	// Get rule info
	info := rules.GetRuleInfoByCode(code)
	// Check if found
	if info == nil {
		fmt.Fprintf(os.Stderr, "Rule not found: %s\n", code)
		OsExit(1)
	}

	// Load example if requested
	if !opts.NoExamples {
		loaded := rules.LoadGoodExamples([]rules.RuleInfo{*info})
		// Check if loaded
		if len(loaded) > 0 {
			info = &loaded[0]
		}
	}

	// Delegate to formatter
	formatter.DisplayRuleDetails(*info)
}
