// Package cmd implements the CLI commands for ktn-linter.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kodflow/ktn-linter/pkg/rules"
)

// jsonRulesFormatter formats rules output as JSON.
type jsonRulesFormatter struct{}

// DisplayCategories shows categories in JSON format.
//
// Params:
//   - categories: list of category names
func (f *jsonRulesFormatter) DisplayCategories(categories []string) {
	var catInfos []categoryInfoJSON
	// Iterate categories
	for _, cat := range categories {
		catRules := rules.GetRuleInfosByCategory(cat)
		catInfos = append(catInfos, categoryInfoJSON{Name: cat, Count: len(catRules)})
	}
	// Encode JSON
	f.encodeCategoriesOutput(categoriesOutputJSON{Categories: catInfos})
}

// encodeCategoriesOutput encodes categories output to JSON.
//
// Params:
//   - output: categories output to encode
func (f *jsonRulesFormatter) encodeCategoriesOutput(output categoriesOutputJSON) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	// Handle encoding error
	if err := encoder.Encode(output); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		OsExit(1)
	}
}

// DisplayCategoryRules shows rules in JSON format.
//
// Params:
//   - category: category name
//   - catRules: list of rules
func (f *jsonRulesFormatter) DisplayCategoryRules(category string, catRules []rules.RuleInfo) {
	// Build and encode category rules
	output := categoryRulesOutputJSON{
		Category: category,
		Rules:    catRules,
	}
	f.encodeCategoryRulesOutput(output)
}

// encodeCategoryRulesOutput encodes category rules output to JSON.
//
// Params:
//   - output: category rules output to encode
func (f *jsonRulesFormatter) encodeCategoryRulesOutput(output categoryRulesOutputJSON) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	// Handle encoding error
	if err := encoder.Encode(output); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		OsExit(1)
	}
}

// DisplayRuleDetails shows rule details in JSON format.
//
// Params:
//   - info: rule information
func (f *jsonRulesFormatter) DisplayRuleDetails(info rules.RuleInfo) {
	// Encode rule info
	f.encodeRuleInfoOutput(info)
}

// encodeRuleInfoOutput encodes rule info to JSON.
//
// Params:
//   - info: rule information to encode
func (f *jsonRulesFormatter) encodeRuleInfoOutput(info rules.RuleInfo) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	// Handle encoding error
	if err := encoder.Encode(info); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		OsExit(1)
	}
}
