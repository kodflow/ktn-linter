package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
)

// TestAnalyzers tests the functionality of Analyzers with error cases
func TestAnalyzers(t *testing.T) {
	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "returns non-empty list",
			check: func(t *testing.T) {
				analyzers := ktntest.Analyzers()
				// Vérification liste non vide
				if len(analyzers) == 0 {
					t.Error("Analyzers should return at least one analyzer")
				}
			},
		},
		{
			name: "all analyzers are valid - error handling",
			check: func(t *testing.T) {
				analyzers := ktntest.Analyzers()
				// Vérification chaque analyzer
				for _, analyzer := range analyzers {
					// Vérification non-nil
					if analyzer == nil {
						t.Error("Analyzers returned nil analyzer")
					}
					// Vérification nom
					if analyzer.Name == "" {
						t.Error("Analyzer has empty name")
					}
					// Vérification doc
					if analyzer.Doc == "" {
						t.Error("Analyzer has empty documentation")
					}
					// Vérification fonction Run
					if analyzer.Run == nil {
						t.Errorf("Analyzer %s has nil Run function", analyzer.Name)
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
