// External tests for the JSON rules formatter.
package cmd_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd"
	"github.com/kodflow/ktn-linter/pkg/rules"
)

// TestDisplayCategories tests DisplayCategories method for JSON formatter.
func TestJSONDisplayCategories(t *testing.T) {
	tests := []struct {
		name         string
		categories   []string
		expectedText string
	}{
		{
			name:         "displays categories as json",
			categories:   []string{"func", "var"},
			expectedText: "categories",
		},
		{
			name:         "displays category name",
			categories:   []string{"func"},
			expectedText: "func",
		},
		{
			name:         "includes count field",
			categories:   []string{"var"},
			expectedText: "count",
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
			formatter := cmd.NewRulesFormatter("json")
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

// TestDisplayCategoryRules tests DisplayCategoryRules method for JSON formatter.
func TestJSONDisplayCategoryRules(t *testing.T) {
	tests := []struct {
		name         string
		category     string
		rules        []rules.RuleInfo
		expectedText string
	}{
		{
			name:     "displays category field",
			category: "func",
			rules: []rules.RuleInfo{
				{Code: "KTN-FUNC-001", Description: "Test rule"},
			},
			expectedText: "category",
		},
		{
			name:     "displays rule code",
			category: "var",
			rules: []rules.RuleInfo{
				{Code: "KTN-VAR-001", Description: "Var rule"},
			},
			expectedText: "KTN-VAR-001",
		},
		{
			name:     "displays rules array",
			category: "func",
			rules: []rules.RuleInfo{
				{Code: "KTN-FUNC-001", Description: "Test"},
			},
			expectedText: "rules",
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
			formatter := cmd.NewRulesFormatter("json")
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

// TestDisplayRuleDetails tests DisplayRuleDetails method for JSON formatter.
func TestJSONDisplayRuleDetails(t *testing.T) {
	tests := []struct {
		name         string
		info         rules.RuleInfo
		expectedText string
	}{
		{
			name: "displays rule code",
			info: rules.RuleInfo{
				Code:        "KTN-FUNC-001",
				Category:    "func",
				Description: "Test description",
			},
			expectedText: "KTN-FUNC-001",
		},
		{
			name: "displays category",
			info: rules.RuleInfo{
				Code:        "KTN-VAR-001",
				Category:    "var",
				Description: "Var rule",
			},
			expectedText: "var",
		},
		{
			name: "displays description",
			info: rules.RuleInfo{
				Code:        "KTN-FUNC-001",
				Category:    "func",
				Description: "Test description for rule",
			},
			expectedText: "Test description for rule",
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
			formatter := cmd.NewRulesFormatter("json")
			formatter.DisplayRuleDetails(tt.info)

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
