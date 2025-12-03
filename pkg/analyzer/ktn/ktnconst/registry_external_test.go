package ktnconst_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst"
)

// TestGetAnalyzers tests GetAnalyzers returns all analyzers
func TestGetAnalyzers(t *testing.T) {
	const EXPECTED_COUNT int = 4
	expectedNames := []string{
		"ktnconst001",
		"ktnconst002",
		"ktnconst003",
		"ktnconst004",
	}

	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "returns expected count",
			check: func(t *testing.T) {
				analyzers := ktnconst.GetAnalyzers()
				// Vérification nombre
				if len(analyzers) != EXPECTED_COUNT {
					t.Errorf("GetAnalyzers() returned %d analyzers, expected %d", len(analyzers), EXPECTED_COUNT)
				}
			},
		},
		{
			name: "all analyzers are non-nil",
			check: func(t *testing.T) {
				analyzers := ktnconst.GetAnalyzers()
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
			name: "analyzers have expected names error cases",
			check: func(t *testing.T) {
				analyzers := ktnconst.GetAnalyzers()
				// Vérification noms
				for i, analyzer := range analyzers {
					// Vérification nom
					if analyzer.Name != expectedNames[i] {
						t.Errorf("Analyzer at index %d has name %q, expected %q", i, analyzer.Name, expectedNames[i])
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
