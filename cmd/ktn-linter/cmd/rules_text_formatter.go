// Package cmd implements the CLI commands for ktn-linter.
package cmd

import (
	"fmt"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/rules"
)

// textSeparatorWidth is the width of the separator line for text output.
const textSeparatorWidth int = 20

// textRulesFormatter formats rules output as plain text.
type textRulesFormatter struct{}

// DisplayCategories shows categories in text format.
//
// Params:
//   - categories: list of category names
func (f *textRulesFormatter) DisplayCategories(categories []string) {
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

// DisplayCategoryRules shows rules in text format.
//
// Params:
//   - category: category name
//   - catRules: list of rules in the category
func (f *textRulesFormatter) DisplayCategoryRules(category string, catRules []rules.RuleInfo) {
	fmt.Printf("KTN-%s Rules\n", strings.ToUpper(category))
	fmt.Println(strings.Repeat("=", textSeparatorWidth))
	fmt.Println()
	// Iterate rules
	for _, rule := range catRules {
		fmt.Printf("  %s: %s\n", rule.Code, rule.Description)
	}
	fmt.Println()
	fmt.Printf("Usage: ktn-linter rules %s <number> for details\n", category)
}

// DisplayRuleDetails shows rule details in text format.
//
// Params:
//   - info: rule information
func (f *textRulesFormatter) DisplayRuleDetails(info rules.RuleInfo) {
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
