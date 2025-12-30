// Package prompt_test provides black-box tests for prompt generation.
package prompt_test

import (
	"bytes"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/orchestrator"
	"github.com/kodflow/ktn-linter/pkg/prompt"
)

// TestNewGenerator tests generator creation.
//
// Params:
//   - t: testing object
func TestNewGenerator(t *testing.T) {
	tests := []struct {
		name    string
		verbose bool
		wantNil bool
	}{
		{
			name:    "creates generator",
			verbose: false,
			wantNil: false,
		},
		{
			name:    "creates verbose generator",
			verbose: true,
			wantNil: false,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create stderr buffer
			var stderr bytes.Buffer

			// Create generator
			gen := prompt.NewGenerator(&stderr, tt.verbose)

			// Verify not nil
			if (gen == nil) != tt.wantNil {
				t.Errorf("NewGenerator() nil = %v, want %v", gen == nil, tt.wantNil)
			}
		})
	}
}

// TestGenerator_Generate tests the Generate public method.
//
// Params:
//   - t: testing object
func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		options  orchestrator.Options
		wantErr  bool
	}{
		{
			name:     "empty patterns",
			patterns: []string{},
			options:  orchestrator.Options{},
			wantErr:  false,
		},
		{
			name:     "valid pattern",
			patterns: []string{"github.com/kodflow/ktn-linter/pkg/prompt"},
			options:  orchestrator.Options{},
			wantErr:  false,
		},
		{
			name:     "invalid pattern",
			patterns: []string{"invalid/package/path/does/not/exist"},
			options:  orchestrator.Options{},
			wantErr:  true,
		},
		{
			name:     "invalid analyzer",
			patterns: []string{"github.com/kodflow/ktn-linter/pkg/prompt"},
			options:  orchestrator.Options{OnlyRule: "INVALID-RULE-999"},
			wantErr:  true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create generator
			var stderr bytes.Buffer
			gen := prompt.NewGenerator(&stderr, false)

			// Generate
			output, err := gen.Generate(tt.patterns, tt.options)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Successful runs must return non-nil output
			if !tt.wantErr && output == nil {
				t.Error("Generate() returned nil output")
			}
		})
	}
}

// TestPromptOutput_Structure tests PromptOutput structure.
//
// Params:
//   - t: testing object
func TestPromptOutput_Structure(t *testing.T) {
	tests := []struct {
		name           string
		output         *prompt.PromptOutput
		wantViolations int
		wantPhases     int
		wantRules      int
	}{
		{
			name: "sample output",
			output: &prompt.PromptOutput{
				TotalViolations: 10,
				TotalRules:      3,
				Phases: []prompt.PhaseGroup{
					{
						Phase: prompt.PhaseStructural,
						Name:  "Structural",
						Rules: []prompt.RuleViolations{
							{Code: "KTN-STRUCT-004", Violations: []prompt.Violation{{}}},
						},
					},
				},
			},
			wantViolations: 10,
			wantPhases:     1,
			wantRules:      1,
		},
		{
			name: "empty output",
			output: &prompt.PromptOutput{
				TotalViolations: 0,
				TotalRules:      0,
				Phases:          []prompt.PhaseGroup{},
			},
			wantViolations: 0,
			wantPhases:     0,
			wantRules:      0,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify structure
			if tt.output.TotalViolations != tt.wantViolations {
				t.Errorf("TotalViolations = %d, want %d", tt.output.TotalViolations, tt.wantViolations)
			}

			// Verify phases
			if len(tt.output.Phases) != tt.wantPhases {
				t.Errorf("len(Phases) = %d, want %d", len(tt.output.Phases), tt.wantPhases)
			}

			// Verify phase rules if phases exist
			if tt.wantPhases > 0 && len(tt.output.Phases[0].Rules) != tt.wantRules {
				t.Errorf("len(Phases[0].Rules) = %d, want %d", len(tt.output.Phases[0].Rules), tt.wantRules)
			}
		})
	}
}
