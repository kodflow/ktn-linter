// Package cmd implements the CLI commands for ktn-linter.
package cmd

import (
	"fmt"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/rules"
)

// markdownRulesFormatter formats rules output as markdown.
type markdownRulesFormatter struct{}

// DisplayCategories shows categories in markdown format.
//
// Params:
//   - categories: list of category names
func (f *markdownRulesFormatter) DisplayCategories(categories []string) {
	fmt.Println("# KTN-Linter Categories")
	fmt.Println()
	// Iterate categories
	for _, cat := range categories {
		// Get rule count for category
		catRules := rules.GetRuleInfosByCategory(cat)
		fmt.Printf("- **%s** (%d rules)\n", cat, len(catRules))
	}
}

// DisplayCategoryRules shows rules in markdown format.
//
// Params:
//   - category: category name
//   - catRules: list of rules
func (f *markdownRulesFormatter) DisplayCategoryRules(category string, catRules []rules.RuleInfo) {
	fmt.Printf("# KTN-%s Rules\n\n", strings.ToUpper(category))
	// Iterate rules
	for _, rule := range catRules {
		fmt.Printf("- **%s**: %s\n", rule.Code, rule.Description)
	}
}

// DisplayRuleDetails shows rule details in markdown format.
//
// Params:
//   - info: rule information
func (f *markdownRulesFormatter) DisplayRuleDetails(info rules.RuleInfo) {
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
