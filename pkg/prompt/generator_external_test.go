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
	// Create stderr buffer
	var stderr bytes.Buffer

	// Create generator
	gen := prompt.NewGenerator(&stderr, false)

	// Verify not nil
	if gen == nil {
		t.Error("NewGenerator() returned nil")
	}
}

// TestGenerator_Generate tests the Generate public method.
//
// Params:
//   - t: testing object
func TestGenerator_Generate(t *testing.T) {
	// Test with empty patterns
	t.Run("empty patterns", func(t *testing.T) {
		// Create generator
		var stderr bytes.Buffer
		gen := prompt.NewGenerator(&stderr, false)

		// Generate with empty patterns (loads current dir)
		output, err := gen.Generate([]string{}, orchestrator.Options{})

		// Should not return error
		if err != nil {
			t.Errorf("Generate() error = %v, want nil", err)
			return
		}

		// Verify output exists
		if output == nil {
			t.Error("Generate() returned nil output")
		}
	})

	// Test with valid pattern
	t.Run("valid pattern", func(t *testing.T) {
		// Create generator
		var stderr bytes.Buffer
		gen := prompt.NewGenerator(&stderr, false)

		// Generate with valid pattern (this package)
		output, err := gen.Generate([]string{"github.com/kodflow/ktn-linter/pkg/prompt"}, orchestrator.Options{})

		// Should not return error
		if err != nil {
			t.Errorf("Generate() error = %v", err)
			return
		}

		// Verify output
		if output == nil {
			t.Error("Generate() returned nil output")
		}
	})

	// Test with invalid pattern
	t.Run("invalid pattern", func(t *testing.T) {
		// Create generator
		var stderr bytes.Buffer
		gen := prompt.NewGenerator(&stderr, false)

		// Generate with invalid pattern
		_, err := gen.Generate([]string{"invalid/package/path/does/not/exist"}, orchestrator.Options{})

		// Should return error from LoadPackages
		if err == nil {
			t.Error("Generate() with invalid pattern should return error")
		}
	})

	// Test with invalid analyzer
	t.Run("invalid analyzer", func(t *testing.T) {
		// Create generator
		var stderr bytes.Buffer
		gen := prompt.NewGenerator(&stderr, false)

		// Generate with invalid analyzer option
		_, err := gen.Generate(
			[]string{"github.com/kodflow/ktn-linter/pkg/prompt"},
			orchestrator.Options{
				OnlyRule: "INVALID-RULE-999",
			},
		)

		// Should return error from SelectAnalyzers
		if err == nil {
			t.Error("Generate() with invalid analyzer should return error")
		}
	})
}

// TestPromptOutput_Structure tests PromptOutput structure.
//
// Params:
//   - t: testing object
func TestPromptOutput_Structure(t *testing.T) {
	// Create a sample output
	output := &prompt.PromptOutput{
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
	}

	// Verify structure
	if output.TotalViolations != 10 {
		t.Errorf("TotalViolations = %d, want 10", output.TotalViolations)
	}

	// Verify phases
	if len(output.Phases) != 1 {
		t.Errorf("len(Phases) = %d, want 1", len(output.Phases))
	}

	// Verify phase rules
	if len(output.Phases[0].Rules) != 1 {
		t.Errorf("len(Phases[0].Rules) = %d, want 1", len(output.Phases[0].Rules))
	}
}

