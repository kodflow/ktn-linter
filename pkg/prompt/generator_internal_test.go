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
	// Define test cases for runLinter
	tests := []struct {
		name        string
		patterns    []string
		opts        orchestrator.Options
		expectError bool
	}{
		{
			name:        "valid pattern returns diagnostics slice",
			patterns:    []string{"github.com/kodflow/ktn-linter/pkg/prompt"},
			opts:        orchestrator.Options{},
			expectError: false,
		},
		{
			name:        "invalid pattern returns error",
			patterns:    []string{"invalid/nonexistent/package/path"},
			opts:        orchestrator.Options{},
			expectError: true,
		},
		{
			name:        "invalid analyzer returns error",
			patterns:    []string{"github.com/kodflow/ktn-linter/pkg/prompt"},
			opts:        orchestrator.Options{OnlyRule: "INVALID-RULE-999"},
			expectError: true,
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var stderr bytes.Buffer
			gen := NewGenerator(&stderr, false)

			// Run linter
			_, err := gen.runLinter(tt.patterns, tt.opts)

			// Verify error expectation symmetrically
			if (err != nil) != tt.expectError {
				t.Fatalf("runLinter() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

// Test_buildModernizeCode tests modernize analyzer code generation.
//
// Params:
//   - t: testing object
func Test_buildModernizeCode(t *testing.T) {
	tests := []struct {
		name         string
		analyzerName string
		want         string
	}{
		{
			name:         "stringscut analyzer",
			analyzerName: "stringscut",
			want:         "KTN-MODERNIZE-001",
		},
		{
			name:         "minmax analyzer",
			analyzerName: "minmax",
			want:         "KTN-MODERNIZE-004",
		},
		{
			name:         "unknown analyzer",
			analyzerName: "unknownanalyzer",
			want:         "",
		},
		{
			name:         "empty analyzer name",
			analyzerName: "",
			want:         "",
		},
		{
			name:         "ktn analyzer (not modernize)",
			analyzerName: "ktnfunc001",
			want:         "",
		},
	}

	// Run test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got := buildModernizeCode(tt.analyzerName)
			// Verify result
			if got != tt.want {
				t.Errorf("buildModernizeCode(%q) = %q, want %q", tt.analyzerName, got, tt.want)
			}
		})
	}
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
		tt := tt // Capture range variable
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
		tt := tt // Capture range variable
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
	// Define test cases for collectViolations
	tests := []struct {
		name         string
		expectedLen  int
	}{
		{
			name:        "empty diagnostics returns empty map",
			expectedLen: 0,
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			gen := &Generator{}
			result := gen.collectViolations(nil)
			// Verify empty result
			if len(result) != tt.expectedLen {
				t.Errorf("collectViolations([]) returned %d rules, want %d", len(result), tt.expectedLen)
			}
		})
	}
}

// Test_Generator_enrichWithMetadata tests metadata enrichment.
//
// Params:
//   - t: testing object
func Test_Generator_enrichWithMetadata(t *testing.T) {
	// Define test cases for enrichWithMetadata
	tests := []struct {
		name         string
		violations   map[string]*RuleViolations
		expectedLen  int
		expectedCode string
	}{
		{
			name:         "empty violations returns empty slice",
			violations:   make(map[string]*RuleViolations),
			expectedLen:  0,
			expectedCode: "",
		},
		{
			name: "enriches violations with metadata",
			violations: map[string]*RuleViolations{
				"KTN-FUNC-001": {
					Code:       "KTN-FUNC-001",
					Violations: []Violation{{FilePath: "test.go", Line: 10}},
				},
			},
			expectedLen:  1,
			expectedCode: "KTN-FUNC-001",
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			gen := &Generator{}
			result := gen.enrichWithMetadata(tt.violations)
			// Verify result length
			if len(result) != tt.expectedLen {
				t.Errorf("enrichWithMetadata() returned %d rules, want %d", len(result), tt.expectedLen)
				return
			}
			// Verify code if expected
			if tt.expectedCode != "" && len(result) > 0 {
				if result[0].Code != tt.expectedCode {
					t.Errorf("enrichWithMetadata() code = %q, want %q", result[0].Code, tt.expectedCode)
				}
			}
		})
	}
}

// Test_Generator_buildOutput tests output construction.
//
// Params:
//   - t: testing object
func Test_Generator_buildOutput(t *testing.T) {
	// Define test cases for buildOutput
	tests := []struct {
		name               string
		rules              []RuleViolations
		phases             []PhaseGroup
		expectedViolations int
		expectedRules      int
		expectedPhases     int
	}{
		{
			name:               "empty output returns zeros",
			rules:              []RuleViolations{},
			phases:             []PhaseGroup{},
			expectedViolations: 0,
			expectedRules:      0,
			expectedPhases:     0,
		},
		{
			name: "with violations counts correctly",
			rules: []RuleViolations{
				{Code: "KTN-FUNC-001", Violations: []Violation{{}, {}}},
				{Code: "KTN-FUNC-002", Violations: []Violation{{}}},
			},
			phases:             []PhaseGroup{{Phase: PhaseLocal}},
			expectedViolations: 3,
			expectedRules:      2,
			expectedPhases:     1,
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			gen := &Generator{}
			output := gen.buildOutput(tt.rules, tt.phases)

			// Verify total violations count
			if output.TotalViolations != tt.expectedViolations {
				t.Errorf("buildOutput() TotalViolations = %d, want %d", output.TotalViolations, tt.expectedViolations)
			}

			// Verify total rules count
			if output.TotalRules != tt.expectedRules {
				t.Errorf("buildOutput() TotalRules = %d, want %d", output.TotalRules, tt.expectedRules)
			}

			// Verify phases count
			if len(output.Phases) != tt.expectedPhases {
				t.Errorf("buildOutput() len(Phases) = %d, want %d", len(output.Phases), tt.expectedPhases)
			}
		})
	}
}
