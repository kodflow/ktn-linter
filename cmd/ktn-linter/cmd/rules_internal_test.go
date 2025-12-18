package cmd

import (
	"bytes"
	"io"
	"os"
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
	tests := []struct {
		name       string
		output     rules.RulesOutput
		useWriter  func() io.Writer
		expectExit bool
		exitCode   int
	}{
		{
			name: "successful JSON encoding",
			output: rules.RulesOutput{
				TotalCount: 1,
				Categories: []string{"func"},
				Rules: []rules.RuleInfo{
					{Code: "KTN-FUNC-001", Category: "func"},
				},
			},
			useWriter: func() io.Writer {
				return &bytes.Buffer{}
			},
			expectExit: false,
			exitCode:   0,
		},
		{
			name: "JSON encoding with multiple rules",
			output: rules.RulesOutput{
				TotalCount: 2,
				Categories: []string{"func", "var"},
				Rules: []rules.RuleInfo{
					{Code: "KTN-FUNC-001", Category: "func", Description: "Test func"},
					{Code: "KTN-VAR-001", Category: "var", Description: "Test var"},
				},
			},
			useWriter: func() io.Writer {
				return &bytes.Buffer{}
			},
			expectExit: false,
			exitCode:   0,
		},
		{
			name: "encoding error exits with code 1",
			output: rules.RulesOutput{
				TotalCount: 1,
				Categories: []string{"test"},
				Rules:      []rules.RuleInfo{{Code: "KTN-TEST-001"}},
			},
			useWriter: func() io.Writer {
				return &failingWriter{}
			},
			expectExit: true,
			exitCode:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			writer := tt.useWriter()

			// Capture stderr for error messages
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			exitCode, didExit := catchExitInCmd(t, func() {
				formatRulesJSON(writer, tt.output)
			})

			w.Close()
			var stderr bytes.Buffer
			stderr.ReadFrom(r)
			os.Stderr = oldStderr

			// Verify exit expectation
			if tt.expectExit && !didExit {
				t.Error("expected exit but did not exit")
			}

			// Verify no exit expectation
			if !tt.expectExit && didExit {
				t.Errorf("unexpected exit with code %d", exitCode)
			}

			// Verify exit code
			if tt.expectExit && exitCode != tt.exitCode {
				t.Errorf("expected exit code %d, got %d", tt.exitCode, exitCode)
			}

			// Check for error message on exit
			if tt.expectExit && !strings.Contains(stderr.String(), "Error encoding JSON") {
				t.Errorf("expected error message in stderr, got: %s", stderr.String())
			}

			// Verify successful output
			if !tt.expectExit {
				if buf, ok := writer.(*bytes.Buffer); ok {
					result := buf.String()
					if !strings.Contains(result, "TotalCount") {
						t.Error("JSON output should contain TotalCount")
					}
					if len(tt.output.Rules) > 0 && !strings.Contains(result, tt.output.Rules[0].Code) {
						t.Errorf("JSON output should contain rule code %s", tt.output.Rules[0].Code)
					}
				}
			}
		})
	}
}

// failingWriter is a writer that always returns an error.
type failingWriter struct{}

// Write implements io.Writer but always returns an error.
//
// Params:
//   - p: bytes to write
//
// Returns:
//   - int: 0 (no bytes written)
//   - error: always returns an error
func (fw *failingWriter) Write(p []byte) (int, error) {
	return 0, os.ErrInvalid
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
