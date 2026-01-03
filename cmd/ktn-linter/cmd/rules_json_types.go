// Package cmd implements the CLI commands for ktn-linter.
package cmd

import "github.com/kodflow/ktn-linter/pkg/rules"

// categoryInfoJSON represents category information for JSON output.
type categoryInfoJSON struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// categoriesOutputJSON represents the JSON output for categories list.
type categoriesOutputJSON struct {
	Categories []categoryInfoJSON `json:"categories"`
}

// categoryRulesOutputJSON represents the JSON output for category rules.
type categoryRulesOutputJSON struct {
	Category string           `json:"category"`
	Rules    []rules.RuleInfo `json:"rules"`
}
