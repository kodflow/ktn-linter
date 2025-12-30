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
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags
			if err := rulesCmd.Flags().Set(flagRulesFormat, tt.wantFormat); err != nil {
				t.Fatalf("failed to set %s: %v", flagRulesFormat, err)
			}
			if err := rulesCmd.Flags().Set(flagRulesNoExamples, "false"); err != nil {
				t.Fatalf("failed to set %s: %v", flagRulesNoExamples, err)
			}

			opts := parseRulesOptions(rulesCmd)
			if opts.Format != tt.wantFormat {
				t.Errorf("parseRulesOptions().Format = %q, want %q", opts.Format, tt.wantFormat)
			}
		})
	}
}

func Test_buildRulesOutput(t *testing.T) {
	// Define test cases for buildRulesOutput
	tests := []struct {
		name          string
		infos         []rules.RuleInfo
		expectedCount int
	}{
		{
			name: "builds output with correct count and rules",
			infos: []rules.RuleInfo{
				{Code: "KTN-FUNC-001", Category: "func", Description: "Test"},
				{Code: "KTN-VAR-001", Category: "var", Description: "Test2"},
			},
			expectedCount: 2,
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			output := buildRulesOutput(tt.infos)

			// Check total count
			if output.TotalCount != tt.expectedCount {
				t.Errorf("buildRulesOutput().TotalCount = %d, want %d", output.TotalCount, tt.expectedCount)
			}

			// Check rules length
			if len(output.Rules) != tt.expectedCount {
				t.Errorf("buildRulesOutput().Rules length = %d, want %d", len(output.Rules), tt.expectedCount)
			}
		})
	}
}

func Test_formatRulesMarkdown(t *testing.T) {
	// Define test cases for markdown formatting
	tests := []struct {
		name           string
		output         rules.RulesOutput
		expectCode     string
		expectSection  string
	}{
		{
			name: "formats rules as markdown with code and examples",
			output: rules.RulesOutput{
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
			},
			expectCode:    "KTN-FUNC-001",
			expectSection: "Good Example",
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			formatRulesMarkdown(&buf, tt.output)

			result := buf.String()
			// Check rule code
			if !strings.Contains(result, tt.expectCode) {
				t.Error("Markdown output should contain rule code")
			}
			// Check section
			if !strings.Contains(result, tt.expectSection) {
				t.Error("Markdown output should contain Good Example section")
			}
		})
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
		tt := tt // Capture range variable
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
	// Define test cases for text formatting
	tests := []struct {
		name           string
		output         rules.RulesOutput
		expectCode     string
		expectCategory string
	}{
		{
			name: "formats rules as text with code and category",
			output: rules.RulesOutput{
				TotalCount: 1,
				Categories: []string{"func"},
				Rules: []rules.RuleInfo{
					{
						Code:        "KTN-FUNC-001",
						Category:    "func",
						Description: "Test description",
					},
				},
			},
			expectCode:     "KTN-FUNC-001",
			expectCategory: "FUNC",
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			formatRulesText(&buf, tt.output)

			result := buf.String()
			// Check rule code
			if !strings.Contains(result, tt.expectCode) {
				t.Error("Text output should contain rule code")
			}
			// Check category
			if !strings.Contains(result, tt.expectCategory) {
				t.Error("Text output should contain category header")
			}
		})
	}
}

func Test_writeRuleMarkdown(t *testing.T) {
	// Define test cases for markdown rule writing
	tests := []struct {
		name           string
		rule           rules.RuleInfo
		expectHeader   string
		expectCategory string
	}{
		{
			name: "writes rule with header and category in markdown",
			rule: rules.RuleInfo{
				Code:        "KTN-TEST-001",
				Category:    "test",
				Description: "Test rule",
				GoodExample: "func Example() {}",
			},
			expectHeader:   "## KTN-TEST-001",
			expectCategory: "**Category**: test",
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writeRuleMarkdown(&buf, tt.rule)

			result := buf.String()
			// Check header
			if !strings.Contains(result, tt.expectHeader) {
				t.Error("Should contain rule header")
			}
			// Check category
			if !strings.Contains(result, tt.expectCategory) {
				t.Error("Should contain category")
			}
		})
	}
}

func Test_writeRuleText(t *testing.T) {
	// Define test cases for text rule writing
	tests := []struct {
		name          string
		rule          rules.RuleInfo
		expectCode    string
		expectSection string
	}{
		{
			name: "writes rule with code and example in text format",
			rule: rules.RuleInfo{
				Code:        "KTN-TEST-001",
				Category:    "test",
				Description: "Test rule",
				GoodExample: "func Example() {}",
			},
			expectCode:    "KTN-TEST-001",
			expectSection: "Good Example",
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writeRuleText(&buf, tt.rule)

			result := buf.String()
			// Check rule code
			if !strings.Contains(result, tt.expectCode) {
				t.Error("Should contain rule code")
			}
			// Check example section
			if !strings.Contains(result, tt.expectSection) {
				t.Error("Should contain example section")
			}
		})
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
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			formatRulesOutput(output, tt.format)
		})
	}
}

// Test_runRules tests the runRules command function.
func Test_runRules(t *testing.T) {
	// Define test cases for runRules
	tests := []struct {
		name       string
		format     string
		noExamples string
	}{
		{name: "executes without panic with default flags", format: "markdown", noExamples: "true"},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Snapshot existing flags to avoid leaking state
			prevFormat, _ := rulesCmd.Flags().GetString(flagRulesFormat)
			prevNoExamples, _ := rulesCmd.Flags().GetString(flagRulesNoExamples)
			t.Cleanup(func() {
				// Best effort cleanup - ignore errors
				_ = rulesCmd.Flags().Set(flagRulesFormat, prevFormat)
				_ = rulesCmd.Flags().Set(flagRulesNoExamples, prevNoExamples)
			})

			// Set flags
			if err := rulesCmd.Flags().Set(flagRulesFormat, tt.format); err != nil {
				t.Fatalf("failed to set %s: %v", flagRulesFormat, err)
			}
			if err := rulesCmd.Flags().Set(flagRulesNoExamples, tt.noExamples); err != nil {
				t.Fatalf("failed to set %s: %v", flagRulesNoExamples, err)
			}

			// runRules writes to stdout, just verify it doesn't panic
			runRules(rulesCmd, []string{})
		})
	}
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
		tt := tt // Capture range variable
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
