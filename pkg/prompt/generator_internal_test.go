// Package prompt provides white-box tests for internal generator functions.
package prompt

import (
	"bytes"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/orchestrator"
)

// Test_Generator_runLinter tests the runLinter private method.
//
// Params:
//   - t: testing object
func Test_Generator_runLinter(t *testing.T) {
	// Test with valid pattern
	t.Run("valid pattern returns diagnostics", func(t *testing.T) {
		var stderr bytes.Buffer
		gen := NewGenerator(&stderr, false)

		// Run linter on self
		diags, err := gen.runLinter(
			[]string{"github.com/kodflow/ktn-linter/pkg/prompt"},
			orchestrator.Options{},
		)

		// Should not error
		if err != nil {
			t.Errorf("runLinter() error = %v", err)
			return
		}

		// Diagnostics should be a slice (possibly empty)
		if diags == nil {
			t.Error("runLinter() returned nil diagnostics")
		}
	})

	// Test with invalid pattern
	t.Run("invalid pattern returns error", func(t *testing.T) {
		var stderr bytes.Buffer
		gen := NewGenerator(&stderr, false)

		// Run linter with invalid pattern
		_, err := gen.runLinter(
			[]string{"invalid/nonexistent/package/path"},
			orchestrator.Options{},
		)

		// Should return error
		if err == nil {
			t.Error("runLinter() with invalid pattern should return error")
		}
	})

	// Test with invalid analyzer option
	t.Run("invalid analyzer returns error", func(t *testing.T) {
		var stderr bytes.Buffer
		gen := NewGenerator(&stderr, false)

		// Run linter with invalid analyzer
		_, err := gen.runLinter(
			[]string{"github.com/kodflow/ktn-linter/pkg/prompt"},
			orchestrator.Options{OnlyRule: "INVALID-RULE-999"},
		)

		// Should return error
		if err == nil {
			t.Error("runLinter() with invalid analyzer should return error")
		}
	})
}

// Test_extractRuleCode tests rule code extraction from various message formats.
//
// Params:
//   - t: testing object
func Test_extractRuleCode(t *testing.T) {
	// Test cases for different message formats
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{
			name:    "bracketed format",
			message: "[KTN-FUNC-001] function too long",
			want:    "KTN-FUNC-001",
		},
		{
			name:    "colon format",
			message: "KTN-FUNC-002: max parameters exceeded",
			want:    "KTN-FUNC-002",
		},
		{
			name:    "space format",
			message: "KTN-CONST-001 constant naming violation",
			want:    "KTN-CONST-001",
		},
		{
			name:    "no closing bracket",
			message: "[KTN-FUNC-003 missing bracket",
			want:    "",
		},
		{
			name:    "no colon or space after code",
			message: "KTN-FUNC-004",
			want:    "",
		},
		{
			name:    "no KTN prefix",
			message: "some other error message",
			want:    "",
		},
		{
			name:    "empty message",
			message: "",
			want:    "",
		},
		{
			name:    "KTN prefix but not rule code",
			message: "KTN-something-wrong",
			want:    "",
		},
		{
			name:    "bracketed at start only",
			message: "[KTN-TEST-001]",
			want:    "KTN-TEST-001",
		},
		{
			name:    "colon at end",
			message: "KTN-VAR-001:",
			want:    "KTN-VAR-001",
		},
	}

	// Run all test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractRuleCode(tt.message)
			// Verify extracted code matches expected
			if got != tt.want {
				t.Errorf("extractRuleCode(%q) = %q, want %q", tt.message, got, tt.want)
			}
		})
	}
}

// Test_extractMessage tests message extraction from various diagnostic formats.
//
// Params:
//   - t: testing object
func Test_extractMessage(t *testing.T) {
	// Test cases for different message formats
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{
			name:    "bracketed format with message",
			message: "[KTN-FUNC-001] function exceeds 35 lines",
			want:    "function exceeds 35 lines",
		},
		{
			name:    "colon format with message",
			message: "KTN-FUNC-002: too many parameters",
			want:    "too many parameters",
		},
		{
			name:    "bracketed without message",
			message: "[KTN-TEST-001]",
			want:    "[KTN-TEST-001]",
		},
		{
			name:    "colon without message",
			message: "KTN-VAR-001:",
			want:    "KTN-VAR-001:",
		},
		{
			name:    "no bracket closing",
			message: "[KTN-FUNC-003 incomplete bracket",
			want:    "[KTN-FUNC-003 incomplete bracket",
		},
		{
			name:    "no KTN prefix",
			message: "regular error message",
			want:    "regular error message",
		},
		{
			name:    "empty message",
			message: "",
			want:    "",
		},
		{
			name:    "KTN prefix without colon",
			message: "KTN-something wrong",
			want:    "KTN-something wrong",
		},
		{
			name:    "bracketed with leading spaces",
			message: "[KTN-CONST-001]   constant must use SCREAMING_SNAKE_CASE",
			want:    "constant must use SCREAMING_SNAKE_CASE",
		},
		{
			name:    "colon with leading spaces",
			message: "KTN-STRUCT-001:   struct naming violation",
			want:    "struct naming violation",
		},
	}

	// Run all test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractMessage(tt.message)
			// Verify extracted message matches expected
			if got != tt.want {
				t.Errorf("extractMessage(%q) = %q, want %q", tt.message, got, tt.want)
			}
		})
	}
}

// Test_Generator_collectViolations tests diagnostic grouping by rule code.
//
// Params:
//   - t: testing object
func Test_Generator_collectViolations(t *testing.T) {
	// Create generator for testing
	gen := &Generator{}

	// Test with empty diagnostics
	t.Run("empty diagnostics", func(t *testing.T) {
		result := gen.collectViolations(nil)
		// Verify empty result
		if len(result) != 0 {
			t.Errorf("collectViolations([]) returned %d rules, want 0", len(result))
		}
	})
}

// Test_Generator_enrichWithMetadata tests metadata enrichment.
//
// Params:
//   - t: testing object
func Test_Generator_enrichWithMetadata(t *testing.T) {
	// Create generator for testing
	gen := &Generator{}

	// Test with empty violations
	t.Run("empty violations", func(t *testing.T) {
		violations := make(map[string]*RuleViolations)
		result := gen.enrichWithMetadata(violations)
		// Verify empty result
		if len(result) != 0 {
			t.Errorf("enrichWithMetadata({}) returned %d rules, want 0", len(result))
		}
	})

	// Test with rule violations
	t.Run("enrich violations", func(t *testing.T) {
		violations := map[string]*RuleViolations{
			"KTN-FUNC-001": {
				Code:       "KTN-FUNC-001",
				Violations: []Violation{{FilePath: "test.go", Line: 10}},
			},
		}
		result := gen.enrichWithMetadata(violations)
		// Verify result length
		if len(result) != 1 {
			t.Errorf("enrichWithMetadata() returned %d rules, want 1", len(result))
			return
		}
		// Verify enrichment occurred (metadata should be populated)
		if result[0].Code != "KTN-FUNC-001" {
			t.Errorf("enrichWithMetadata() code = %q, want KTN-FUNC-001", result[0].Code)
		}
	})
}

// Test_Generator_buildOutput tests output construction.
//
// Params:
//   - t: testing object
func Test_Generator_buildOutput(t *testing.T) {
	// Create generator for testing
	gen := &Generator{}

	// Test with empty rules and phases
	t.Run("empty output", func(t *testing.T) {
		output := gen.buildOutput([]RuleViolations{}, []PhaseGroup{})
		// Verify zero violations
		if output.TotalViolations != 0 {
			t.Errorf("buildOutput() TotalViolations = %d, want 0", output.TotalViolations)
		}
		// Verify zero rules
		if output.TotalRules != 0 {
			t.Errorf("buildOutput() TotalRules = %d, want 0", output.TotalRules)
		}
	})

	// Test with rules and violations
	t.Run("with violations", func(t *testing.T) {
		rules := []RuleViolations{
			{Code: "KTN-FUNC-001", Violations: []Violation{{}, {}}},
			{Code: "KTN-FUNC-002", Violations: []Violation{{}}},
		}
		phases := []PhaseGroup{{Phase: PhaseLocal}}

		output := gen.buildOutput(rules, phases)

		// Verify total violations count
		if output.TotalViolations != 3 {
			t.Errorf("buildOutput() TotalViolations = %d, want 3", output.TotalViolations)
		}

		// Verify total rules count
		if output.TotalRules != 2 {
			t.Errorf("buildOutput() TotalRules = %d, want 2", output.TotalRules)
		}

		// Verify phases
		if len(output.Phases) != 1 {
			t.Errorf("buildOutput() len(Phases) = %d, want 1", len(output.Phases))
		}
	})
}
