package ktnfunc_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
)

// TestGetAnalyzers tests GetAnalyzers returns all analyzers
func TestGetAnalyzers(t *testing.T) {
	const MIN_EXPECTED_COUNT int = 14
	expectedNames := map[string]bool{
		"ktnfunc001": true, // Max 35 lines
		"ktnfunc002": true, // Max 5 parameters
		"ktnfunc003": true, // No magic numbers
		"ktnfunc004": true, // No naked returns
		"ktnfunc005": true, // Max cyclomatic complexity 10
		"ktnfunc006": true, // Error last
		"ktnfunc007": true, // Documentation stricte
		"ktnfunc008": true, // Context must be first parameter
		"ktnfunc009": true, // No side effects in getters
		"ktnfunc010": true, // Named returns for >3 return values
		"ktnfunc011": true, // Comments on branches/returns
		"ktnfunc012": true, // No else after return/continue/break
		"ktnfunc013": true, // Unused parameters must be marked
		"ktnfunc014": true, // Private functions must be used in production
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
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t)
		})
	}
}
