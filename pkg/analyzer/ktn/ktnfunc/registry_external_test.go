package ktnfunc_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
)

// TestGetAnalyzers tests GetAnalyzers returns all analyzers
func TestGetAnalyzers(t *testing.T) {
	const MIN_EXPECTED_COUNT int = 12 // FUNC-007 et FUNC-009 déplacés vers COMMENT
	expectedNames := map[string]bool{
		"ktnfunc001": true, // Error last
		"ktnfunc002": true, // Context first
		"ktnfunc003": true, // No else after return
		"ktnfunc004": true, // Private functions used
		"ktnfunc005": true, // Max 35 lines
		"ktnfunc006": true, // Max 5 parameters
		"ktnfunc007": true, // No side effects in getters (ex-008)
		"ktnfunc008": true, // Unused parameters prefixed (ex-010)
		"ktnfunc009": true, // No magic numbers (ex-011)
		"ktnfunc010": true, // No naked returns (ex-012)
		"ktnfunc011": true, // Max cyclomatic complexity (ex-013)
		"ktnfunc012": true, // Named returns >3 values (ex-014)
	}

	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "returns minimum expected count",
			check: func(t *testing.T) {
				analyzers := ktnfunc.GetAnalyzers()
				// Vérification nombre minimum
				if len(analyzers) < MIN_EXPECTED_COUNT {
					t.Errorf("GetAnalyzers() returned %d analyzers, expected at least %d", len(analyzers), MIN_EXPECTED_COUNT)
				}
			},
		},
		{
			name: "all analyzers are non-nil",
			check: func(t *testing.T) {
				analyzers := ktnfunc.GetAnalyzers()
				// Vérification non-nil
				for i, analyzer := range analyzers {
					// Vérification analyzer
					if analyzer == nil {
						t.Errorf("Analyzer at index %d is nil", i)
					}
				}
			},
		},
		{
			name: "all analyzer names are expected",
			check: func(t *testing.T) {
				analyzers := ktnfunc.GetAnalyzers()
				// Vérification noms attendus
				for _, analyzer := range analyzers {
					// Vérification nom dans map
					if !expectedNames[analyzer.Name] {
						t.Errorf("Unexpected analyzer name: %s", analyzer.Name)
					}
				}
			},
		},
	}

	// Exécution tests
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t)
		})
	}
}
