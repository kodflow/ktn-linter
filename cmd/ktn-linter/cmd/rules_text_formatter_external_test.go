// External tests for the text rules formatter.
package cmd_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd"
	"github.com/kodflow/ktn-linter/pkg/rules"
)

// TestTextDisplayCategories tests DisplayCategories method for text formatter.
func TestTextDisplayCategories(t *testing.T) {
	tests := []struct {
		name         string
		categories   []string
		expectedText string
	}{
		{
			name:         "displays header",
			categories:   []string{"func", "var"},
			expectedText: "KTN-Linter Categories",
		},
		{
			name:         "displays separator",
			categories:   []string{"func"},
			expectedText: "=====================",
		},
		{
			name:         "displays usage hint",
			categories:   []string{"var"},
			expectedText: "Usage: ktn-linter rules",
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
			formatter := cmd.NewRulesFormatter("text")
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

// TestTextDisplayCategoryRules tests DisplayCategoryRules method for text formatter.
func TestTextDisplayCategoryRules(t *testing.T) {
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
			expectedText: "KTN-FUNC Rules",
		},
		{
			name:     "displays separator line",
			category: "var",
			rules: []rules.RuleInfo{
				{Code: "KTN-VAR-001", Description: "Test description"},
			},
			expectedText: "====================",
		},
		{
			name:     "displays rule code",
			category: "func",
			rules: []rules.RuleInfo{
				{Code: "KTN-FUNC-001", Description: "My description"},
			},
			expectedText: "KTN-FUNC-001:",
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
			formatter := cmd.NewRulesFormatter("text")
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

// TestTextDisplayRuleDetails tests DisplayRuleDetails method for text formatter.
func TestTextDisplayRuleDetails(t *testing.T) {
	tests := []struct {
		name         string
		rule         rules.RuleInfo
		expectedText string
	}{
		{
			name: "displays rule code",
			rule: rules.RuleInfo{
				Code:        "KTN-FUNC-001",
				Category:    "func",
				Description: "Function description",
			},
			expectedText: "KTN-FUNC-001",
		},
		{
			name: "displays category label",
			rule: rules.RuleInfo{
				Code:        "KTN-VAR-001",
				Category:    "var",
				Description: "Variable description",
			},
			expectedText: "Category: var",
		},
		{
			name: "displays good example header when present",
			rule: rules.RuleInfo{
				Code:        "KTN-FUNC-001",
				Category:    "func",
				Description: "Description",
				GoodExample: "func good() {}",
			},
			expectedText: "Good Example:",
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
			formatter := cmd.NewRulesFormatter("text")
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
