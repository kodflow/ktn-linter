// Package cmd implements the CLI commands for ktn-linter.
package cmd

import (
	"encoding/json"
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
	Args: cobra.MaximumNArgs(2),
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

// runRules executes the rules command.
//
// Params:
//   - cmd: Cobra command (used to get flags)
//   - args: command line arguments (category and/or rule number)
//
// Returns: none
func runRules(cmd *cobra.Command, args []string) {
	// Parse command-specific flags
	opts := parseRulesOptions(cmd)

	// Handle hierarchical navigation based on args
	switch len(args) {
	// No args: list categories
	case 0:
		displayCategories(opts.Format)
	// One arg: category name or full rule code
	case 1:
		handleSingleArg(args[0], opts)
	// Two args: category and rule number
	case 2:
		handleCategoryAndRule(args[0], args[1], opts)
	}
}

// displayCategories shows all available categories.
//
// Params:
//   - format: output format (text, markdown, json)
func displayCategories(format string) {
	categories := rules.GetCategories()

	// Select format
	switch strings.ToLower(format) {
	// JSON format
	case "json":
		displayCategoriesJSON(categories)
	// Markdown format
	case "markdown", "md":
		displayCategoriesMarkdown(categories)
	// Text format (default)
	default:
		displayCategoriesText(categories)
	}
}

// displayCategoriesText shows categories in text format.
//
// Params:
//   - categories: list of category names
func displayCategoriesText(categories []string) {
	fmt.Println("KTN-Linter Categories")
	fmt.Println("=====================")
	fmt.Println()
	// Iterate categories
	for _, cat := range categories {
		// Get rule count for category
		catRules := rules.GetRuleInfosByCategory(cat)
		fmt.Printf("  %s (%d rules)\n", cat, len(catRules))
	}
	fmt.Println()
	fmt.Println("Usage: ktn-linter rules <category> to see rules")
}

// displayCategoriesMarkdown shows categories in markdown format.
//
// Params:
//   - categories: list of category names
func displayCategoriesMarkdown(categories []string) {
	fmt.Println("# KTN-Linter Categories")
	fmt.Println()
	// Iterate categories
	for _, cat := range categories {
		// Get rule count for category
		catRules := rules.GetRuleInfosByCategory(cat)
		fmt.Printf("- **%s** (%d rules)\n", cat, len(catRules))
	}
}

// displayCategoriesJSON shows categories in JSON format.
//
// Params:
//   - categories: list of category names
func displayCategoriesJSON(categories []string) {
	// Build JSON structure
	type categoryInfo struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	var catInfos []categoryInfo
	// Iterate categories
	for _, cat := range categories {
		catRules := rules.GetRuleInfosByCategory(cat)
		catInfos = append(catInfos, categoryInfo{Name: cat, Count: len(catRules)})
	}
	// Encode JSON
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	// Handle encoding error
	if err := encoder.Encode(map[string]any{"categories": catInfos}); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		OsExit(1)
	}
}

// handleSingleArg handles one argument (category or full rule code).
//
// Params:
//   - arg: category name or rule code (e.g., "func" or "KTN-FUNC-001")
//   - opts: rules options
func handleSingleArg(arg string, opts rulesOptions) {
	// Check if it's a full rule code (KTN-XXX-YYY)
	if strings.HasPrefix(strings.ToUpper(arg), "KTN-") {
		// Display single rule details
		displayRuleDetails(strings.ToUpper(arg), opts)
		return
	}

	// Otherwise treat as category
	displayCategoryRules(arg, opts.Format)
}

// displayCategoryRules shows all rules in a category.
//
// Params:
//   - category: category name
//   - format: output format
func displayCategoryRules(category string, format string) {
	catRules := rules.GetRuleInfosByCategory(category)

	// Check if category exists
	if len(catRules) == 0 {
		fmt.Fprintf(os.Stderr, "Category not found: %s\n", category)
		fmt.Fprintf(os.Stderr, "Available categories: %s\n", strings.Join(rules.GetCategories(), ", "))
		OsExit(1)
	}

	// Select format
	switch strings.ToLower(format) {
	// JSON format
	case "json":
		displayCategoryRulesJSON(category, catRules)
	// Markdown format
	case "markdown", "md":
		displayCategoryRulesMarkdown(category, catRules)
	// Text format (default)
	default:
		displayCategoryRulesText(category, catRules)
	}
}

// displayCategoryRulesText shows rules in text format.
//
// Params:
//   - category: category name
//   - catRules: list of rules in the category
func displayCategoryRulesText(category string, catRules []rules.RuleInfo) {
	fmt.Printf("KTN-%s Rules\n", strings.ToUpper(category))
	fmt.Println(strings.Repeat("=", 20))
	fmt.Println()
	// Iterate rules
	for _, rule := range catRules {
		fmt.Printf("  %s: %s\n", rule.Code, rule.Description)
	}
	fmt.Println()
	fmt.Printf("Usage: ktn-linter rules %s <number> for details\n", category)
}

// displayCategoryRulesMarkdown shows rules in markdown format.
//
// Params:
//   - category: category name
//   - catRules: list of rules
func displayCategoryRulesMarkdown(category string, catRules []rules.RuleInfo) {
	fmt.Printf("# KTN-%s Rules\n\n", strings.ToUpper(category))
	// Iterate rules
	for _, rule := range catRules {
		fmt.Printf("- **%s**: %s\n", rule.Code, rule.Description)
	}
}

// displayCategoryRulesJSON shows rules in JSON format.
//
// Params:
//   - category: category name
//   - catRules: list of rules
func displayCategoryRulesJSON(category string, catRules []rules.RuleInfo) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	// Handle encoding error
	if err := encoder.Encode(map[string]any{
		"category": category,
		"rules":    catRules,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		OsExit(1)
	}
}

// handleCategoryAndRule handles two arguments (category and rule number).
//
// Params:
//   - category: category name (e.g., "func")
//   - ruleNum: rule number (e.g., "001")
//   - opts: rules options
func handleCategoryAndRule(category, ruleNum string, opts rulesOptions) {
	// Validate rule number format (must be 3 digits)
	if !isValidRuleNumber(ruleNum) {
		fmt.Fprintf(os.Stderr, "Invalid rule number format: %s (expected 3 digits, e.g., 001)\n", ruleNum)
		OsExit(1)
	}

	// Build full rule code
	code := fmt.Sprintf("KTN-%s-%s", strings.ToUpper(category), ruleNum)
	displayRuleDetails(code, opts)
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
	// Check length
	if len(ruleNum) != 3 {
		// Invalid length
		return false
	}

	// Check all characters are digits
	for _, c := range ruleNum {
		if c < '0' || c > '9' {
			// Non-digit character found
			return false
		}
	}

	// Valid format
	return true
}

// displayRuleDetails shows detailed information for a single rule.
//
// Params:
//   - code: full rule code (e.g., "KTN-FUNC-001")
//   - opts: rules options
func displayRuleDetails(code string, opts rulesOptions) {
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

	// Select format
	switch strings.ToLower(opts.Format) {
	// JSON format
	case "json":
		displayRuleDetailsJSON(*info)
	// Markdown format
	case "markdown", "md":
		displayRuleDetailsMarkdown(*info)
	// Text format (default)
	default:
		displayRuleDetailsText(*info)
	}
}

// displayRuleDetailsText shows rule details in text format.
//
// Params:
//   - info: rule information
func displayRuleDetailsText(info rules.RuleInfo) {
	fmt.Printf("%s\n", info.Code)
	fmt.Println(strings.Repeat("=", len(info.Code)))
	fmt.Println()
	fmt.Printf("Category: %s\n", info.Category)
	fmt.Printf("Description: %s\n", info.Description)
	// Show example if available
	if info.GoodExample != "" {
		fmt.Println()
		fmt.Println("Good Example:")
		fmt.Println("-------------")
		// Iterate lines
		for line := range strings.SplitSeq(info.GoodExample, "\n") {
			fmt.Printf("  %s\n", line)
		}
	}
}

// displayRuleDetailsMarkdown shows rule details in markdown format.
//
// Params:
//   - info: rule information
func displayRuleDetailsMarkdown(info rules.RuleInfo) {
	fmt.Printf("# %s\n\n", info.Code)
	fmt.Printf("**Category**: %s\n\n", info.Category)
	fmt.Printf("%s\n\n", info.Description)
	// Show example if available
	if info.GoodExample != "" {
		fmt.Println("## Good Example")
		fmt.Println()
		fmt.Println("```go")
		fmt.Print(info.GoodExample)
		fmt.Println("```")
	}
}

// displayRuleDetailsJSON shows rule details in JSON format.
//
// Params:
//   - info: rule information
func displayRuleDetailsJSON(info rules.RuleInfo) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	// Handle encoding error
	if err := encoder.Encode(info); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		OsExit(1)
	}
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

