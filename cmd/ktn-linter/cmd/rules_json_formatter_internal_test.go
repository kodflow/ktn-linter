// Internal tests for the JSON rules formatter.
package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/rules"
)

// Test_jsonRulesFormatter_encodeCategoriesOutput tests the encodeCategoriesOutput method.
func Test_jsonRulesFormatter_encodeCategoriesOutput(t *testing.T) {
	tests := []struct {
		name           string
		output         categoriesOutputJSON
		expectedOutput string
	}{
		{
			name: "encodes categories correctly",
			output: categoriesOutputJSON{
				Categories: []categoryInfoJSON{
					{Name: "func", Count: 10},
					{Name: "var", Count: 5},
				},
			},
			expectedOutput: "func",
		},
		{
			name: "encodes empty categories",
			output: categoriesOutputJSON{
				Categories: []categoryInfoJSON{},
			},
			expectedOutput: "categories",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Execute encoding
			f := &jsonRulesFormatter{}
			f.encodeCategoriesOutput(tt.output)

			w.Close()
			var stdout bytes.Buffer
			stdout.ReadFrom(r)
			os.Stdout = oldStdout

			// Verify output contains expected
			if !strings.Contains(stdout.String(), tt.expectedOutput) {
				t.Errorf("output should contain %q, got %s", tt.expectedOutput, stdout.String())
			}
		})
	}
}

// Test_jsonRulesFormatter_encodeCategoryRulesOutput tests the encodeCategoryRulesOutput method.
func Test_jsonRulesFormatter_encodeCategoryRulesOutput(t *testing.T) {
	tests := []struct {
		name           string
		output         categoryRulesOutputJSON
		expectedOutput string
	}{
		{
			name: "encodes category rules correctly",
			output: categoryRulesOutputJSON{
				Category: "func",
				Rules: []rules.RuleInfo{
					{Code: "KTN-FUNC-001", Description: "Test rule"},
				},
			},
			expectedOutput: "KTN-FUNC-001",
		},
		{
			name: "encodes empty rules",
			output: categoryRulesOutputJSON{
				Category: "var",
				Rules:    []rules.RuleInfo{},
			},
			expectedOutput: "var",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Execute encoding
			f := &jsonRulesFormatter{}
			f.encodeCategoryRulesOutput(tt.output)

			w.Close()
			var stdout bytes.Buffer
			stdout.ReadFrom(r)
			os.Stdout = oldStdout

			// Verify output contains expected
			if !strings.Contains(stdout.String(), tt.expectedOutput) {
				t.Errorf("output should contain %q, got %s", tt.expectedOutput, stdout.String())
			}
		})
	}
}

// Test_jsonRulesFormatter_encodeRuleInfoOutput tests the encodeRuleInfoOutput method.
func Test_jsonRulesFormatter_encodeRuleInfoOutput(t *testing.T) {
	tests := []struct {
		name           string
		info           rules.RuleInfo
		expectedOutput string
	}{
		{
			name: "encodes rule info correctly",
			info: rules.RuleInfo{
				Code:        "KTN-FUNC-001",
				Category:    "func",
				Description: "Test description",
			},
			expectedOutput: "KTN-FUNC-001",
		},
		{
			name: "encodes rule with example",
			info: rules.RuleInfo{
				Code:        "KTN-VAR-001",
				Category:    "var",
				Description: "Test var rule",
				GoodExample: "var x int",
			},
			expectedOutput: "GoodExample",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Execute encoding
			f := &jsonRulesFormatter{}
			f.encodeRuleInfoOutput(tt.info)

			w.Close()
			var stdout bytes.Buffer
			stdout.ReadFrom(r)
			os.Stdout = oldStdout

			// Verify output contains expected
			if !strings.Contains(stdout.String(), tt.expectedOutput) {
				t.Errorf("output should contain %q, got %s", tt.expectedOutput, stdout.String())
			}
		})
	}
}
