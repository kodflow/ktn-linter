package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func Test_parseRulesOptions(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		noExamples bool
		wantFormat string
	}{
		{
			name:       "default format is text",
			format:     "text",
			noExamples: false,
			wantFormat: "text",
		},
		{
			name:       "markdown format",
			format:     "markdown",
			noExamples: false,
			wantFormat: "markdown",
		},
		{
			name:       "json format",
			format:     "json",
			noExamples: true,
			wantFormat: "json",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Set flags
			if err := rulesCmd.Flags().Set(flagRulesFormat, tt.format); err != nil {
				t.Fatalf("failed to set %s: %v", flagRulesFormat, err)
			}
			noExamplesStr := "false"
			if tt.noExamples {
				noExamplesStr = "true"
			}
			if err := rulesCmd.Flags().Set(flagRulesNoExamples, noExamplesStr); err != nil {
				t.Fatalf("failed to set %s: %v", flagRulesNoExamples, err)
			}

			opts := parseRulesOptions(rulesCmd)
			// Verify format
			if opts.Format != tt.wantFormat {
				t.Errorf("parseRulesOptions().Format = %q, want %q", opts.Format, tt.wantFormat)
			}
			// Verify noExamples
			if opts.NoExamples != tt.noExamples {
				t.Errorf("parseRulesOptions().NoExamples = %v, want %v", opts.NoExamples, tt.noExamples)
			}
		})
	}
}

func Test_NewRulesFormatter(t *testing.T) {
	tests := []struct {
		name         string
		format       string
		expectedType string
	}{
		{
			name:         "text format",
			format:       "text",
			expectedType: "*cmd.textRulesFormatter",
		},
		{
			name:         "markdown format",
			format:       "markdown",
			expectedType: "*cmd.markdownRulesFormatter",
		},
		{
			name:         "md alias for markdown",
			format:       "md",
			expectedType: "*cmd.markdownRulesFormatter",
		},
		{
			name:         "json format",
			format:       "json",
			expectedType: "*cmd.jsonRulesFormatter",
		},
		{
			name:         "unknown format defaults to text",
			format:       "unknown",
			expectedType: "*cmd.textRulesFormatter",
		},
		{
			name:         "empty format defaults to text",
			format:       "",
			expectedType: "*cmd.textRulesFormatter",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewRulesFormatter(tt.format)
			// Verify formatter is not nil
			if formatter == nil {
				t.Fatal("NewRulesFormatter returned nil")
			}
			// Verify type by checking interface implementation
			switch tt.expectedType {
			case "*cmd.textRulesFormatter":
				if _, ok := formatter.(*textRulesFormatter); !ok {
					t.Errorf("expected *textRulesFormatter, got %T", formatter)
				}
			case "*cmd.markdownRulesFormatter":
				if _, ok := formatter.(*markdownRulesFormatter); !ok {
					t.Errorf("expected *markdownRulesFormatter, got %T", formatter)
				}
			case "*cmd.jsonRulesFormatter":
				if _, ok := formatter.(*jsonRulesFormatter); !ok {
					t.Errorf("expected *jsonRulesFormatter, got %T", formatter)
				}
			}
		})
	}
}

func Test_displayCategories(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		expectText string
	}{
		{
			name:       "text format shows categories",
			format:     "text",
			expectText: "Categories",
		},
		{
			name:       "markdown format shows categories",
			format:     "markdown",
			expectText: "# KTN-Linter Categories",
		},
		{
			name:       "json format produces valid output",
			format:     "json",
			expectText: "categories",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Create formatter and call function
			formatter := NewRulesFormatter(tt.format)
			displayCategories(formatter)

			w.Close()
			var stdout bytes.Buffer
			stdout.ReadFrom(r)
			os.Stdout = oldStdout

			result := stdout.String()
			// Verify expected text
			if !strings.Contains(result, tt.expectText) {
				t.Errorf("displayCategories(%q) output should contain %q, got: %s",
					tt.format, tt.expectText, result)
			}
		})
	}
}

func Test_displayCategoryRules(t *testing.T) {
	tests := []struct {
		name       string
		category   string
		format     string
		expectText string
		expectExit bool
		exitCode   int
	}{
		{
			name:       "valid category text format",
			category:   "func",
			format:     "text",
			expectText: "KTN-FUNC",
			expectExit: false,
		},
		{
			name:       "valid category markdown format",
			category:   "func",
			format:     "markdown",
			expectText: "# KTN-FUNC",
			expectExit: false,
		},
		{
			name:       "valid category json format",
			category:   "func",
			format:     "json",
			expectText: "category",
			expectExit: false,
		},
		{
			name:       "invalid category exits with error",
			category:   "nonexistent",
			format:     "text",
			expectText: "",
			expectExit: true,
			exitCode:   1,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Create formatter
			formatter := NewRulesFormatter(tt.format)

			exitCode, didExit := catchExitInCmd(t, func() {
				displayCategoryRules(tt.category, formatter)
			})

			w.Close()
			var stdout bytes.Buffer
			stdout.ReadFrom(r)
			os.Stdout = oldStdout

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
			// Verify expected text if not expecting exit
			if !tt.expectExit && tt.expectText != "" {
				result := stdout.String()
				if !strings.Contains(result, tt.expectText) {
					t.Errorf("output should contain %q, got: %s", tt.expectText, result)
				}
			}
		})
	}
}

func Test_displayRuleDetails(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		format     string
		noExamples bool
		expectText string
		expectExit bool
		exitCode   int
	}{
		{
			name:       "valid rule text format",
			code:       "KTN-FUNC-001",
			format:     "text",
			noExamples: true,
			expectText: "KTN-FUNC-001",
			expectExit: false,
		},
		{
			name:       "valid rule markdown format",
			code:       "KTN-FUNC-001",
			format:     "markdown",
			noExamples: true,
			expectText: "# KTN-FUNC-001",
			expectExit: false,
		},
		{
			name:       "valid rule json format",
			code:       "KTN-FUNC-001",
			format:     "json",
			noExamples: true,
			expectText: "KTN-FUNC-001",
			expectExit: false,
		},
		{
			name:       "invalid rule exits with error",
			code:       "KTN-INVALID-999",
			format:     "text",
			noExamples: true,
			expectText: "",
			expectExit: true,
			exitCode:   1,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Create formatter
			formatter := NewRulesFormatter(tt.format)
			opts := rulesOptions{
				Format:     tt.format,
				NoExamples: tt.noExamples,
			}

			exitCode, didExit := catchExitInCmd(t, func() {
				displayRuleDetails(tt.code, formatter, opts)
			})

			w.Close()
			var stdout bytes.Buffer
			stdout.ReadFrom(r)
			os.Stdout = oldStdout

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
			// Verify expected text if not expecting exit
			if !tt.expectExit && tt.expectText != "" {
				result := stdout.String()
				if !strings.Contains(result, tt.expectText) {
					t.Errorf("output should contain %q, got: %s", tt.expectText, result)
				}
			}
		})
	}
}

func Test_handleSingleArg(t *testing.T) {
	tests := []struct {
		name       string
		arg        string
		format     string
		expectText string
		expectExit bool
	}{
		{
			name:       "category name shows rules",
			arg:        "func",
			format:     "text",
			expectText: "KTN-FUNC",
			expectExit: false,
		},
		{
			name:       "full rule code shows details",
			arg:        "KTN-FUNC-001",
			format:     "text",
			expectText: "KTN-FUNC-001",
			expectExit: false,
		},
		{
			name:       "lowercase rule code converted to uppercase",
			arg:        "ktn-func-001",
			format:     "text",
			expectText: "KTN-FUNC-001",
			expectExit: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Create formatter
			formatter := NewRulesFormatter(tt.format)
			opts := rulesOptions{
				Format:     tt.format,
				NoExamples: true,
			}

			_, didExit := catchExitInCmd(t, func() {
				handleSingleArg(tt.arg, formatter, opts)
			})

			w.Close()
			var stdout bytes.Buffer
			stdout.ReadFrom(r)
			os.Stdout = oldStdout

			// Verify exit expectation
			if tt.expectExit != didExit {
				t.Errorf("expected exit=%v, got exit=%v", tt.expectExit, didExit)
			}
			// Verify expected text if not expecting exit
			if !tt.expectExit && tt.expectText != "" {
				result := stdout.String()
				if !strings.Contains(result, tt.expectText) {
					t.Errorf("output should contain %q, got: %s", tt.expectText, result)
				}
			}
		})
	}
}

func Test_handleCategoryAndRule(t *testing.T) {
	tests := []struct {
		name       string
		category   string
		ruleNum    string
		format     string
		expectText string
		expectExit bool
		exitCode   int
	}{
		{
			name:       "valid category and rule number",
			category:   "func",
			ruleNum:    "001",
			format:     "text",
			expectText: "KTN-FUNC-001",
			expectExit: false,
		},
		{
			name:       "invalid rule number exits",
			category:   "func",
			ruleNum:    "999",
			format:     "text",
			expectText: "",
			expectExit: true,
			exitCode:   1,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Create formatter
			formatter := NewRulesFormatter(tt.format)
			opts := rulesOptions{
				Format:     tt.format,
				NoExamples: true,
			}

			exitCode, didExit := catchExitInCmd(t, func() {
				handleCategoryAndRule(tt.category, tt.ruleNum, formatter, opts)
			})

			w.Close()
			var stdout bytes.Buffer
			stdout.ReadFrom(r)
			os.Stdout = oldStdout

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
			// Verify expected text if not expecting exit
			if !tt.expectExit && tt.expectText != "" {
				result := stdout.String()
				if !strings.Contains(result, tt.expectText) {
					t.Errorf("output should contain %q, got: %s", tt.expectText, result)
				}
			}
		})
	}
}

func Test_runRules(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		format string
	}{
		{
			name:   "no args shows categories",
			args:   []string{},
			format: "text",
		},
		{
			name:   "category arg shows rules",
			args:   []string{"func"},
			format: "text",
		},
		{
			name:   "category and number shows details",
			args:   []string{"func", "001"},
			format: "text",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Set flags
			if err := rulesCmd.Flags().Set(flagRulesFormat, tt.format); err != nil {
				t.Fatalf("failed to set format: %v", err)
			}
			if err := rulesCmd.Flags().Set(flagRulesNoExamples, "true"); err != nil {
				t.Fatalf("failed to set no-examples: %v", err)
			}

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Should not panic
			runRules(rulesCmd, tt.args)

			w.Close()
			var stdout bytes.Buffer
			stdout.ReadFrom(r)
			os.Stdout = oldStdout

			// Verify output is not empty
			if stdout.Len() == 0 {
				t.Error("expected non-empty output")
			}
		})
	}
}

func Test_isValidRuleNumber(t *testing.T) {
	tests := []struct {
		name     string
		ruleNum  string
		expected bool
	}{
		{
			name:     "valid three digits",
			ruleNum:  "001",
			expected: true,
		},
		{
			name:     "valid mid range",
			ruleNum:  "123",
			expected: true,
		},
		{
			name:     "valid all nines",
			ruleNum:  "999",
			expected: true,
		},
		{
			name:     "invalid too short",
			ruleNum:  "01",
			expected: false,
		},
		{
			name:     "invalid too long",
			ruleNum:  "0001",
			expected: false,
		},
		{
			name:     "invalid contains letter",
			ruleNum:  "00a",
			expected: false,
		},
		{
			name:     "invalid all letters",
			ruleNum:  "abc",
			expected: false,
		},
		{
			name:     "invalid empty string",
			ruleNum:  "",
			expected: false,
		},
		{
			name:     "invalid special characters",
			ruleNum:  "0-1",
			expected: false,
		},
		{
			name:     "invalid spaces",
			ruleNum:  " 01",
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := isValidRuleNumber(tt.ruleNum)
			// Verify result
			if result != tt.expected {
				t.Errorf("isValidRuleNumber(%q) = %v, want %v", tt.ruleNum, result, tt.expected)
			}
		})
	}
}

func Test_handleCategoryAndRule_InvalidFormat(t *testing.T) {
	tests := []struct {
		name       string
		category   string
		ruleNum    string
		expectExit bool
		exitCode   int
	}{
		{
			name:       "invalid rule number too short",
			category:   "func",
			ruleNum:    "01",
			expectExit: true,
			exitCode:   1,
		},
		{
			name:       "invalid rule number too long",
			category:   "func",
			ruleNum:    "0001",
			expectExit: true,
			exitCode:   1,
		},
		{
			name:       "invalid rule number contains letter",
			category:   "func",
			ruleNum:    "00a",
			expectExit: true,
			exitCode:   1,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Capture stderr for error message
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			// Create formatter
			formatter := NewRulesFormatter("text")
			opts := rulesOptions{
				Format:     "text",
				NoExamples: true,
			}

			exitCode, didExit := catchExitInCmd(t, func() {
				handleCategoryAndRule(tt.category, tt.ruleNum, formatter, opts)
			})

			w.Close()
			var stderr bytes.Buffer
			stderr.ReadFrom(r)
			os.Stderr = oldStderr

			// Verify exit expectation
			if tt.expectExit && !didExit {
				t.Error("expected exit but did not exit")
			}
			// Verify exit code
			if tt.expectExit && exitCode != tt.exitCode {
				t.Errorf("expected exit code %d, got %d", tt.exitCode, exitCode)
			}
			// Verify error message contains expected text
			if tt.expectExit && !strings.Contains(stderr.String(), "Invalid rule number format") {
				t.Errorf("expected error message about invalid format, got: %s", stderr.String())
			}
		})
	}
}

func TestRulesCmdStructure(t *testing.T) {
	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "Use field is correct",
			check: func(t *testing.T) {
				expectedUse := "rules [category] [rule-number]"
				// Verify Use
				if rulesCmd.Use != expectedUse {
					t.Errorf("expected Use=%q, got %q", expectedUse, rulesCmd.Use)
				}
			},
		},
		{
			name: "Short description is not empty",
			check: func(t *testing.T) {
				// Verify Short
				if rulesCmd.Short == "" {
					t.Error("Short description should not be empty")
				}
			},
		},
		{
			name: "Long description contains usage examples",
			check: func(t *testing.T) {
				// Verify Long contains examples
				if !strings.Contains(rulesCmd.Long, "ktn-linter rules") {
					t.Error("Long description should contain usage examples")
				}
			},
		},
		{
			name: "Run function is not nil",
			check: func(t *testing.T) {
				// Verify Run
				if rulesCmd.Run == nil {
					t.Error("Run function should not be nil")
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t)
		})
	}
}
