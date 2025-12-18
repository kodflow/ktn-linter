// Package prompt provides white-box tests for internal formatter functions.
package prompt

import (
	"bytes"
	"strings"
	"testing"
)

// Test_MarkdownFormatter_writeHeader tests the writeHeader private method.
//
// Params:
//   - t: testing object
func Test_MarkdownFormatter_writeHeader(t *testing.T) {
	tests := []struct {
		name       string
		output     *PromptOutput
		wantTitle  bool
		wantSumary string
	}{
		{
			name: "empty output",
			output: &PromptOutput{
				TotalViolations: 0,
				TotalRules:      0,
			},
			wantTitle:  true,
			wantSumary: "0 violations across 0 rules",
		},
		{
			name: "with violations",
			output: &PromptOutput{
				TotalViolations: 5,
				TotalRules:      2,
			},
			wantTitle:  true,
			wantSumary: "5 violations across 2 rules",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			f := &MarkdownFormatter{writer: &buf}

			// Call private function
			f.writeHeader(tt.output)

			result := buf.String()

			// Verify title
			if tt.wantTitle && !strings.Contains(result, "# KTN-Linter Correction Prompt") {
				t.Error("writeHeader() should contain title")
			}

			// Verify summary
			if !strings.Contains(result, tt.wantSumary) {
				t.Errorf("writeHeader() should contain %q", tt.wantSumary)
			}
		})
	}
}

// Test_MarkdownFormatter_writeInstructions tests the writeInstructions private method.
//
// Params:
//   - t: testing object
func Test_MarkdownFormatter_writeInstructions(t *testing.T) {
	tests := []struct {
		name        string
		wantSection bool
		wantContent []string
	}{
		{
			name:        "instructions section",
			wantSection: true,
			wantContent: []string{
				"## Instructions",
				"Ce prompt guide la correction",
				"Phase 1",
			},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			f := &MarkdownFormatter{writer: &buf}

			// Call private function
			f.writeInstructions()

			result := buf.String()

			// Verify section header
			if tt.wantSection && !strings.Contains(result, "## Instructions") {
				t.Error("writeInstructions() should contain section header")
			}

			// Verify content
			for _, content := range tt.wantContent {
				if !strings.Contains(result, content) {
					t.Errorf("writeInstructions() should contain %q", content)
				}
			}
		})
	}
}

// Test_MarkdownFormatter_writePhase tests the writePhase private method.
//
// Params:
//   - t: testing object
func Test_MarkdownFormatter_writePhase(t *testing.T) {
	tests := []struct {
		name       string
		phase      PhaseGroup
		phaseNum   int
		wantHeader string
		wantRerun  bool
	}{
		{
			name: "structural phase with rerun",
			phase: PhaseGroup{
				Phase:       PhaseStructural,
				Name:        "Structural Changes",
				Description: "May create/move/delete files.",
				NeedsRerun:  true,
				Rules: []RuleViolations{
					{Code: "KTN-STRUCT-004", Violations: []Violation{{}}},
				},
			},
			phaseNum:   1,
			wantHeader: "## Phase 1: Structural Changes",
			wantRerun:  true,
		},
		{
			name: "local phase without rerun",
			phase: PhaseGroup{
				Phase:       PhaseLocal,
				Name:        "Local Fixes",
				Description: "Code modifications within existing files.",
				NeedsRerun:  false,
				Rules: []RuleViolations{
					{Code: "KTN-FUNC-001", Violations: []Violation{{}}},
				},
			},
			phaseNum:   2,
			wantHeader: "## Phase 2: Local Fixes",
			wantRerun:  false,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			f := &MarkdownFormatter{writer: &buf}

			// Call private function
			f.writePhase(tt.phase, tt.phaseNum)

			result := buf.String()

			// Verify header
			if !strings.Contains(result, tt.wantHeader) {
				t.Errorf("writePhase() should contain %q", tt.wantHeader)
			}

			// Verify rerun warning
			hasRerun := strings.Contains(result, "Re-executez le linter")
			if hasRerun != tt.wantRerun {
				t.Errorf("writePhase() rerun warning = %v, want %v", hasRerun, tt.wantRerun)
			}
		})
	}
}

// Test_MarkdownFormatter_writeRule tests the writeRule private method.
//
// Params:
//   - t: testing object
func Test_MarkdownFormatter_writeRule(t *testing.T) {
	tests := []struct {
		name        string
		rule        RuleViolations
		wantHeader  string
		wantExample bool
	}{
		{
			name: "rule with example",
			rule: RuleViolations{
				Code:        "KTN-FUNC-001",
				Category:    "func",
				Description: "Function too long",
				GoodExample: "func short() {}",
				Violations:  []Violation{{FilePath: "test.go", Line: 10}},
			},
			wantHeader:  "### KTN-FUNC-001",
			wantExample: true,
		},
		{
			name: "rule without example",
			rule: RuleViolations{
				Code:        "KTN-VAR-002",
				Category:    "var",
				Description: "Variable naming",
				GoodExample: "",
				Violations:  []Violation{{FilePath: "test.go", Line: 5}},
			},
			wantHeader:  "### KTN-VAR-002",
			wantExample: false,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			f := &MarkdownFormatter{writer: &buf}

			// Call private function
			f.writeRule(tt.rule)

			result := buf.String()

			// Verify header
			if !strings.Contains(result, tt.wantHeader) {
				t.Errorf("writeRule() should contain %q", tt.wantHeader)
			}

			// Verify example section
			hasExample := strings.Contains(result, "#### Good Example")
			if hasExample != tt.wantExample {
				t.Errorf("writeRule() example section = %v, want %v", hasExample, tt.wantExample)
			}
		})
	}
}

// Test_MarkdownFormatter_writeGoodExample tests the writeGoodExample private method.
//
// Params:
//   - t: testing object
func Test_MarkdownFormatter_writeGoodExample(t *testing.T) {
	tests := []struct {
		name        string
		example     string
		wantCode    bool
		wantNewline bool
	}{
		{
			name:        "example with newline",
			example:     "func foo() {}\n",
			wantCode:    true,
			wantNewline: false,
		},
		{
			name:        "example without newline",
			example:     "func bar() {}",
			wantCode:    true,
			wantNewline: true,
		},
		{
			name:        "empty example",
			example:     "",
			wantCode:    true,
			wantNewline: false,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			f := &MarkdownFormatter{writer: &buf}

			// Call private function
			f.writeGoodExample(tt.example)

			result := buf.String()

			// Verify code block
			if tt.wantCode && !strings.Contains(result, "```go") {
				t.Error("writeGoodExample() should contain go code block")
			}

			// Verify closing
			if !strings.Contains(result, "```\n") {
				t.Error("writeGoodExample() should have closing code fence")
			}
		})
	}
}

// Test_MarkdownFormatter_writeViolations tests the writeViolations private method.
//
// Params:
//   - t: testing object
func Test_MarkdownFormatter_writeViolations(t *testing.T) {
	tests := []struct {
		name       string
		violations []Violation
		wantCount  string
	}{
		{
			name:       "empty violations",
			violations: []Violation{},
			wantCount:  "0 violations",
		},
		{
			name: "single violation",
			violations: []Violation{
				{FilePath: "test.go", Line: 10},
			},
			wantCount: "1 violations",
		},
		{
			name: "multiple violations",
			violations: []Violation{
				{FilePath: "a.go", Line: 1},
				{FilePath: "b.go", Line: 2},
				{FilePath: "c.go", Line: 3},
			},
			wantCount: "3 violations",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			f := &MarkdownFormatter{writer: &buf}

			// Call private function
			f.writeViolations(tt.violations)

			result := buf.String()

			// Verify count header
			if !strings.Contains(result, tt.wantCount) {
				t.Errorf("writeViolations() should contain %q", tt.wantCount)
			}
		})
	}
}

// Test_MarkdownFormatter_writeViolation tests the writeViolation private method.
//
// Params:
//   - t: testing object
func Test_MarkdownFormatter_writeViolation(t *testing.T) {
	tests := []struct {
		name    string
		v       Violation
		want    string
		notWant string
	}{
		{
			name: "with message",
			v: Violation{
				FilePath: "pkg/test.go",
				Line:     25,
				Message:  "error not last",
			},
			want:    "- [ ] `pkg/test.go:25` - error not last",
			notWant: "",
		},
		{
			name: "without message",
			v: Violation{
				FilePath: "pkg/other.go",
				Line:     10,
				Message:  "",
			},
			want:    "- [ ] `pkg/other.go:10`",
			notWant: "- [ ] `pkg/other.go:10` -",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			f := &MarkdownFormatter{writer: &buf}

			// Call private function
			f.writeViolation(tt.v)

			result := buf.String()

			// Verify format
			if !strings.Contains(result, tt.want) {
				t.Errorf("writeViolation() = %q, want to contain %q", result, tt.want)
			}

			// Verify not contains unwanted
			if tt.notWant != "" && strings.Contains(result, tt.notWant) {
				t.Errorf("writeViolation() = %q, should not contain %q", result, tt.notWant)
			}
		})
	}
}

// TestMarkdownFormatter_Format tests the Format public method.
//
// Params:
//   - t: testing object
func TestMarkdownFormatter_Format(t *testing.T) {
	tests := []struct {
		name   string
		output *PromptOutput
	}{
		{
			name: "complete format",
			output: &PromptOutput{
				TotalViolations: 1,
				TotalRules:      1,
				Phases: []PhaseGroup{
					{
						Phase: PhaseLocal,
						Name:  "Local Fixes",
						Rules: []RuleViolations{
							{
								Code:       "KTN-TEST-001",
								Category:   "test",
								Violations: []Violation{{FilePath: "t.go", Line: 1}},
							},
						},
					},
				},
			},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			f := NewMarkdownFormatter(&buf)

			// Call Format
			f.Format(tt.output)

			// Verify output is not empty
			if buf.Len() == 0 {
				t.Error("Format() produced empty output")
			}
		})
	}
}
