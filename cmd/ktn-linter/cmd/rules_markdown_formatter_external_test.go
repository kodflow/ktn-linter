// External tests for the Markdown rules formatter.
package cmd_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd"
	"github.com/kodflow/ktn-linter/pkg/rules"
)

// TestMarkdownDisplayCategories tests DisplayCategories method for markdown formatter.
func TestMarkdownDisplayCategories(t *testing.T) {
	tests := []struct {
		name         string
		categories   []string
		expectedText string
	}{
		{
			name:         "displays header",
			categories:   []string{"func", "var"},
			expectedText: "# KTN-Linter Categories",
		},
		{
			name:         "displays category name",
			categories:   []string{"func"},
			expectedText: "**func**",
		},
		{
			name:         "displays rules count",
			categories:   []string{"var"},
			expectedText: "rules)",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Create formatter and display
			formatter := cmd.NewRulesFormatter("markdown")
			formatter.DisplayCategories(tt.categories)

			w.Close()
			var stdout bytes.Buffer
			stdout.ReadFrom(r)
			os.Stdout = oldStdout

			// Verify output contains expected text
			if !strings.Contains(stdout.String(), tt.expectedText) {
				t.Errorf("output should contain %q, got %s", tt.expectedText, stdout.String())
			}
		})
	}
}

// TestMarkdownDisplayCategoryRules tests DisplayCategoryRules method for markdown formatter.
func TestMarkdownDisplayCategoryRules(t *testing.T) {
	tests := []struct {
		name         string
		category     string
		rules        []rules.RuleInfo
		expectedText string
	}{
		{
			name:     "displays category header",
			category: "func",
			rules: []rules.RuleInfo{
				{Code: "KTN-FUNC-001", Description: "Test rule"},
			},
			expectedText: "# KTN-FUNC Rules",
		},
		{
			name:     "displays rule code in bold",
			category: "var",
			rules: []rules.RuleInfo{
				{Code: "KTN-VAR-001", Description: "Test description"},
			},
			expectedText: "**KTN-VAR-001**",
		},
		{
			name:     "displays rule description",
			category: "func",
			rules: []rules.RuleInfo{
				{Code: "KTN-FUNC-001", Description: "My test description"},
			},
			expectedText: "My test description",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Create formatter and display
			formatter := cmd.NewRulesFormatter("markdown")
			formatter.DisplayCategoryRules(tt.category, tt.rules)

			w.Close()
			var stdout bytes.Buffer
			stdout.ReadFrom(r)
			os.Stdout = oldStdout

			// Verify output contains expected text
			if !strings.Contains(stdout.String(), tt.expectedText) {
				t.Errorf("output should contain %q, got %s", tt.expectedText, stdout.String())
			}
		})
	}
}

// TestMarkdownDisplayRuleDetails tests DisplayRuleDetails method for markdown formatter.
func TestMarkdownDisplayRuleDetails(t *testing.T) {
	tests := []struct {
		name         string
		rule         rules.RuleInfo
		expectedText string
	}{
		{
			name: "displays rule code header",
			rule: rules.RuleInfo{
				Code:        "KTN-FUNC-001",
				Category:    "func",
				Description: "Function description",
			},
			expectedText: "# KTN-FUNC-001",
		},
		{
			name: "displays category",
			rule: rules.RuleInfo{
				Code:        "KTN-VAR-001",
				Category:    "var",
				Description: "Variable description",
			},
			expectedText: "**Category**: var",
		},
		{
			name: "displays good example when present",
			rule: rules.RuleInfo{
				Code:        "KTN-FUNC-001",
				Category:    "func",
				Description: "Description",
				GoodExample: "func good() {}",
			},
			expectedText: "## Good Example",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Create formatter and display
			formatter := cmd.NewRulesFormatter("markdown")
			formatter.DisplayRuleDetails(tt.rule)

			w.Close()
			var stdout bytes.Buffer
			stdout.ReadFrom(r)
			os.Stdout = oldStdout

			// Verify output contains expected text
			if !strings.Contains(stdout.String(), tt.expectedText) {
				t.Errorf("output should contain %q, got %s", tt.expectedText, stdout.String())
			}
		})
	}
}
