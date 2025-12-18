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
	// Create buffer
	var buf bytes.Buffer

	// Create formatter
	f := prompt.NewMarkdownFormatter(&buf)

	// Verify not nil
	if f == nil {
		t.Error("NewMarkdownFormatter() returned nil")
	}
}

// TestMarkdownFormatter_Format_EmptyOutput tests formatting empty output.
//
// Params:
//   - t: testing object
func TestMarkdownFormatter_Format_EmptyOutput(t *testing.T) {
	// Create buffer and formatter
	var buf bytes.Buffer
	f := prompt.NewMarkdownFormatter(&buf)

	// Format empty output
	output := &prompt.PromptOutput{
		TotalViolations: 0,
		TotalRules:      0,
		Phases:          []prompt.PhaseGroup{},
	}
	f.Format(output)

	// Verify output contains header
	result := buf.String()
	if !strings.Contains(result, "# KTN-Linter Correction Prompt") {
		t.Error("Output should contain header")
	}

	// Verify summary shows 0
	if !strings.Contains(result, "0 violations") {
		t.Error("Output should show 0 violations")
	}
}

// TestMarkdownFormatter_Format_WithViolations tests formatting with violations.
//
// Params:
//   - t: testing object
func TestMarkdownFormatter_Format_WithViolations(t *testing.T) {
	// Create buffer and formatter
	var buf bytes.Buffer
	f := prompt.NewMarkdownFormatter(&buf)

	// Create output with violations
	output := &prompt.PromptOutput{
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
	}
	f.Format(output)

	// Get result
	result := buf.String()

	// Verify header
	if !strings.Contains(result, "# KTN-Linter Correction Prompt") {
		t.Error("Output should contain header")
	}

	// Verify summary
	if !strings.Contains(result, "3 violations across 1 rules") {
		t.Error("Output should show correct summary")
	}

	// Verify phase header
	if !strings.Contains(result, "## Phase 1: Local Fixes") {
		t.Error("Output should contain phase header")
	}

	// Verify rule header
	if !strings.Contains(result, "### KTN-FUNC-001") {
		t.Error("Output should contain rule header")
	}

	// Verify category
	if !strings.Contains(result, "**Category**: func") {
		t.Error("Output should contain category")
	}

	// Verify good example
	if !strings.Contains(result, "#### Good Example") {
		t.Error("Output should contain good example section")
	}

	// Verify violations
	if !strings.Contains(result, "- [ ] `pkg/service/handler.go:23`") {
		t.Error("Output should contain checkbox violations")
	}
}

// TestMarkdownFormatter_Format_MultiplePhases tests formatting multiple phases.
//
// Params:
//   - t: testing object
func TestMarkdownFormatter_Format_MultiplePhases(t *testing.T) {
	// Create buffer and formatter
	var buf bytes.Buffer
	f := prompt.NewMarkdownFormatter(&buf)

	// Create output with multiple phases
	output := &prompt.PromptOutput{
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
	}
	f.Format(output)

	// Get result
	result := buf.String()

	// Verify both phases
	if !strings.Contains(result, "## Phase 1: Structural Changes") {
		t.Error("Output should contain structural phase")
	}
	if !strings.Contains(result, "## Phase 2: Comments & Documentation") {
		t.Error("Output should contain comment phase")
	}

	// Verify rerun warning for structural phase
	if !strings.Contains(result, "Re-executez le linter apres cette phase") {
		t.Error("Output should contain rerun warning")
	}
}

// TestMarkdownFormatter_Format_ViolationWithMessage tests violation message formatting.
//
// Params:
//   - t: testing object
func TestMarkdownFormatter_Format_ViolationWithMessage(t *testing.T) {
	// Create buffer and formatter
	var buf bytes.Buffer
	f := prompt.NewMarkdownFormatter(&buf)

	// Create output with message
	output := &prompt.PromptOutput{
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
	}
	f.Format(output)

	// Verify message is included
	result := buf.String()
	if !strings.Contains(result, "- [ ] `test.go:10` - variable name too long") {
		t.Error("Output should contain violation with message")
	}
}

// TestMarkdownFormatter_Format_ViolationWithoutMessage tests violation without message.
//
// Params:
//   - t: testing object
func TestMarkdownFormatter_Format_ViolationWithoutMessage(t *testing.T) {
	// Create buffer and formatter
	var buf bytes.Buffer
	f := prompt.NewMarkdownFormatter(&buf)

	// Create output without message
	output := &prompt.PromptOutput{
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
	}
	f.Format(output)

	// Verify format without message
	result := buf.String()
	if !strings.Contains(result, "- [ ] `test.go:10`\n") {
		t.Error("Output should contain violation without trailing dash")
	}
}
