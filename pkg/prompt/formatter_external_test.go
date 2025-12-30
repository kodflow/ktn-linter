// Package prompt_test provides black-box tests for markdown formatting.
package prompt_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/prompt"
)

// TestNewMarkdownFormatter tests formatter creation.
//
// Params:
//   - t: testing object
func TestNewMarkdownFormatter(t *testing.T) {
	tests := []struct {
		name    string
		wantNil bool
	}{
		{
			name:    "creates formatter",
			wantNil: false,
		},
	}

	// Run test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create buffer
			var buf bytes.Buffer

			// Create formatter
			f := prompt.NewMarkdownFormatter(&buf)

			// Verify not nil
			if (f == nil) != tt.wantNil {
				t.Errorf("NewMarkdownFormatter() nil = %v, want %v", f == nil, tt.wantNil)
			}
		})
	}
}

// TestMarkdownFormatter_Format tests various format scenarios.
//
// Params:
//   - t: testing object
func TestMarkdownFormatter_Format(t *testing.T) {
	tests := []struct {
		name       string
		output     *prompt.PromptOutput
		wantParts  []string
		notWant    []string
	}{
		{
			name: "empty output",
			output: &prompt.PromptOutput{
				TotalViolations: 0,
				TotalRules:      0,
				Phases:          []prompt.PhaseGroup{},
			},
			wantParts: []string{
				"# KTN-Linter Correction Prompt",
				"0 violations",
			},
			notWant: []string{},
		},
		{
			name: "with violations",
			output: &prompt.PromptOutput{
				TotalViolations: 3,
				TotalRules:      1,
				Phases: []prompt.PhaseGroup{
					{
						Phase:       prompt.PhaseLocal,
						Name:        "Local Fixes",
						Description: "Code modifications within existing files.",
						NeedsRerun:  false,
						Rules: []prompt.RuleViolations{
							{
								Code:        "KTN-FUNC-001",
								Category:    "func",
								Description: "Error must be last return value",
								GoodExample: "func Do() (*Result, error) { return nil, nil }",
								Violations: []prompt.Violation{
									{FilePath: "pkg/service/handler.go", Line: 23, Message: "error not last"},
									{FilePath: "pkg/service/handler.go", Line: 45, Message: "error not last"},
									{FilePath: "pkg/api/routes.go", Line: 12, Message: "error not last"},
								},
							},
						},
					},
				},
			},
			wantParts: []string{
				"# KTN-Linter Correction Prompt",
				"3 violations across 1 rules",
				"## Phase 1: Local Fixes",
				"### KTN-FUNC-001",
				"**Category**: func",
				"#### Good Example",
				"- [ ] `pkg/service/handler.go:23`",
			},
			notWant: []string{},
		},
		{
			name: "multiple phases",
			output: &prompt.PromptOutput{
				TotalViolations: 4,
				TotalRules:      2,
				Phases: []prompt.PhaseGroup{
					{
						Phase:       prompt.PhaseStructural,
						Name:        "Structural Changes",
						Description: "May create/move/delete files.",
						NeedsRerun:  true,
						Rules: []prompt.RuleViolations{
							{
								Code:        "KTN-STRUCT-004",
								Category:    "struct",
								Description: "One struct per file",
								Violations: []prompt.Violation{
									{FilePath: "pkg/models/entities.go", Line: 15},
								},
							},
						},
					},
					{
						Phase:       prompt.PhaseComment,
						Name:        "Comments & Documentation",
						Description: "Apply last after all code is finalized.",
						NeedsRerun:  false,
						Rules: []prompt.RuleViolations{
							{
								Code:        "KTN-COMMENT-005",
								Category:    "comment",
								Description: "Function documentation required",
								Violations: []prompt.Violation{
									{FilePath: "pkg/utils/helpers.go", Line: 10},
									{FilePath: "pkg/utils/helpers.go", Line: 25},
									{FilePath: "pkg/utils/helpers.go", Line: 40},
								},
							},
						},
					},
				},
			},
			wantParts: []string{
				"## Phase 1: Structural Changes",
				"## Phase 2: Comments & Documentation",
				"Re-executez le linter apres cette phase",
			},
			notWant: []string{},
		},
		{
			name: "violation with message",
			output: &prompt.PromptOutput{
				TotalViolations: 1,
				TotalRules:      1,
				Phases: []prompt.PhaseGroup{
					{
						Phase: prompt.PhaseLocal,
						Name:  "Local Fixes",
						Rules: []prompt.RuleViolations{
							{
								Code:     "KTN-VAR-002",
								Category: "var",
								Violations: []prompt.Violation{
									{
										FilePath: "test.go",
										Line:     10,
										Message:  "variable name too long",
									},
								},
							},
						},
					},
				},
			},
			wantParts: []string{
				"- [ ] `test.go:10` - variable name too long",
			},
			notWant: []string{},
		},
		{
			name: "violation without message",
			output: &prompt.PromptOutput{
				TotalViolations: 1,
				TotalRules:      1,
				Phases: []prompt.PhaseGroup{
					{
						Phase: prompt.PhaseLocal,
						Name:  "Local Fixes",
						Rules: []prompt.RuleViolations{
							{
								Code:     "KTN-VAR-002",
								Category: "var",
								Violations: []prompt.Violation{
									{FilePath: "test.go", Line: 10, Message: ""},
								},
							},
						},
					},
				},
			},
			wantParts: []string{
				"- [ ] `test.go:10`\n",
			},
			notWant: []string{
				"- [ ] `test.go:10` -",
			},
		},
		{
			name: "complete format",
			output: &prompt.PromptOutput{
				TotalViolations: 1,
				TotalRules:      1,
				Phases: []prompt.PhaseGroup{
					{
						Phase: prompt.PhaseLocal,
						Name:  "Local Fixes",
						Rules: []prompt.RuleViolations{
							{
								Code:       "KTN-TEST-001",
								Category:   "test",
								Violations: []prompt.Violation{{FilePath: "t.go", Line: 1}},
							},
						},
					},
				},
			},
			wantParts: []string{
				"# KTN-Linter Correction Prompt",
			},
			notWant: []string{},
		},
	}

	// Run test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create buffer and formatter
			var buf bytes.Buffer
			f := prompt.NewMarkdownFormatter(&buf)

			// Format output
			f.Format(tt.output)

			// Get result
			result := buf.String()

			// Verify expected parts
			for _, part := range tt.wantParts {
				if !strings.Contains(result, part) {
					t.Errorf("Format() output should contain %q", part)
				}
			}

			// Verify unwanted parts are absent
			for _, part := range tt.notWant {
				if strings.Contains(result, part) {
					t.Errorf("Format() output should not contain %q", part)
				}
			}

			// Verify output is not empty
			if buf.Len() == 0 {
				t.Error("Format() produced empty output")
			}
		})
	}
}
