package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/rules"
)

func Test_parseRulesOptions(t *testing.T) {
	tests := []struct {
		name       string
		wantFormat string
	}{
		{
			name:       "default format is markdown",
			wantFormat: "markdown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags
			rulesCmd.Flags().Set(flagRulesFormat, tt.wantFormat)
			rulesCmd.Flags().Set(flagRulesNoExamples, "false")

			opts := parseRulesOptions(rulesCmd)
			if opts.Format != tt.wantFormat {
				t.Errorf("parseRulesOptions().Format = %q, want %q", opts.Format, tt.wantFormat)
			}
		})
	}
}

func Test_buildRulesOutput(t *testing.T) {
	infos := []rules.RuleInfo{
		{Code: "KTN-FUNC-001", Category: "func", Description: "Test"},
		{Code: "KTN-VAR-001", Category: "var", Description: "Test2"},
	}

	output := buildRulesOutput(infos)

	if output.TotalCount != 2 {
		t.Errorf("buildRulesOutput().TotalCount = %d, want 2", output.TotalCount)
	}

	if len(output.Rules) != 2 {
		t.Errorf("buildRulesOutput().Rules length = %d, want 2", len(output.Rules))
	}
}

func Test_formatRulesMarkdown(t *testing.T) {
	output := rules.RulesOutput{
		TotalCount: 1,
		Categories: []string{"func"},
		Rules: []rules.RuleInfo{
			{
				Code:        "KTN-FUNC-001",
				Category:    "func",
				Description: "Test description",
				GoodExample: "func Test() {}",
			},
		},
	}

	var buf bytes.Buffer
	formatRulesMarkdown(&buf, output)

	result := buf.String()
	if !strings.Contains(result, "KTN-FUNC-001") {
		t.Error("Markdown output should contain rule code")
	}
	if !strings.Contains(result, "Good Example") {
		t.Error("Markdown output should contain Good Example section")
	}
}

func Test_formatRulesJSON(t *testing.T) {
	output := rules.RulesOutput{
		TotalCount: 1,
		Categories: []string{"func"},
		Rules: []rules.RuleInfo{
			{Code: "KTN-FUNC-001", Category: "func"},
		},
	}

	var buf bytes.Buffer
	formatRulesJSON(&buf, output)

	result := buf.String()
	if !strings.Contains(result, "TotalCount") {
		t.Error("JSON output should contain TotalCount")
	}
	if !strings.Contains(result, "KTN-FUNC-001") {
		t.Error("JSON output should contain rule code")
	}
}

func Test_formatRulesText(t *testing.T) {
	output := rules.RulesOutput{
		TotalCount: 1,
		Categories: []string{"func"},
		Rules: []rules.RuleInfo{
			{
				Code:        "KTN-FUNC-001",
				Category:    "func",
				Description: "Test description",
			},
		},
	}

	var buf bytes.Buffer
	formatRulesText(&buf, output)

	result := buf.String()
	if !strings.Contains(result, "KTN-FUNC-001") {
		t.Error("Text output should contain rule code")
	}
	if !strings.Contains(result, "FUNC") {
		t.Error("Text output should contain category header")
	}
}

func Test_writeRuleMarkdown(t *testing.T) {
	rule := rules.RuleInfo{
		Code:        "KTN-TEST-001",
		Category:    "test",
		Description: "Test rule",
		GoodExample: "func Example() {}",
	}

	var buf bytes.Buffer
	writeRuleMarkdown(&buf, rule)

	result := buf.String()
	if !strings.Contains(result, "## KTN-TEST-001") {
		t.Error("Should contain rule header")
	}
	if !strings.Contains(result, "**Category**: test") {
		t.Error("Should contain category")
	}
}

func Test_writeRuleText(t *testing.T) {
	rule := rules.RuleInfo{
		Code:        "KTN-TEST-001",
		Category:    "test",
		Description: "Test rule",
		GoodExample: "func Example() {}",
	}

	var buf bytes.Buffer
	writeRuleText(&buf, rule)

	result := buf.String()
	if !strings.Contains(result, "KTN-TEST-001") {
		t.Error("Should contain rule code")
	}
	if !strings.Contains(result, "Good Example") {
		t.Error("Should contain example section")
	}
}

func Test_formatRulesOutput(t *testing.T) {
	output := rules.RulesOutput{
		TotalCount: 1,
		Categories: []string{"func"},
		Rules:      []rules.RuleInfo{{Code: "KTN-FUNC-001"}},
	}

	tests := []struct {
		name   string
		format string
	}{
		{name: "markdown format", format: "markdown"},
		{name: "json format", format: "json"},
		{name: "text format", format: "text"},
		{name: "default format", format: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			formatRulesOutput(output, tt.format)
		})
	}
}

// Test_runRules tests the runRules command function.
func Test_runRules(t *testing.T) {
	// Test with default flags - should not panic
	rulesCmd.Flags().Set(flagRulesFormat, "markdown")
	rulesCmd.Flags().Set(flagRulesNoExamples, "true")

	// runRules writes to stdout, just verify it doesn't panic
	runRules(rulesCmd, []string{})
}

// Test_getRulesWithFilters tests the getRulesWithFilters function.
func Test_getRulesWithFilters(t *testing.T) {
	tests := []struct {
		name       string
		category   string
		onlyRule   string
		wantResult bool
	}{
		{
			name:       "no filters returns all rules",
			category:   "",
			onlyRule:   "",
			wantResult: true,
		},
		{
			name:       "filter by category",
			category:   "func",
			onlyRule:   "",
			wantResult: true,
		},
		{
			name:       "filter by single rule",
			category:   "",
			onlyRule:   "KTN-FUNC-001",
			wantResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set flags via rootCmd
			flags := rootCmd.PersistentFlags()
			flags.Set(flagCategory, tt.category)
			flags.Set(flagOnlyRule, tt.onlyRule)

			infos := getRulesWithFilters()

			// Verify result is valid
			if tt.wantResult && infos == nil {
				t.Error("Expected non-nil result")
			}

			// Reset flags
			flags.Set(flagCategory, "")
			flags.Set(flagOnlyRule, "")
		})
	}
}
