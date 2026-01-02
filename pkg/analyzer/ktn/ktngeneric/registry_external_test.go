package ktngeneric_test

import (
	"testing"

	ktngeneric "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktngeneric"
)

func TestAnalyzers(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "returns non-empty list",
			testFunc: func(t *testing.T) {
				analyzers := ktngeneric.Analyzers()
				// Verifier que la liste n'est pas vide
				if len(analyzers) == 0 {
					t.Error("Analyzers() returned empty list")
				}
			},
		},
		{
			name: "all analyzers have name and doc",
			testFunc: func(t *testing.T) {
				analyzers := ktngeneric.Analyzers()
				// Parcourir les analyseurs
				for _, a := range analyzers {
					// Verifier le nom
					if a.Name == "" {
						t.Error("analyzer has empty name")
					}
					// Verifier la documentation
					if a.Doc == "" {
						t.Errorf("analyzer %s has empty doc", a.Name)
					}
				}
			},
		},
		{
			name: "returns expected count",
			testFunc: func(t *testing.T) {
				analyzers := ktngeneric.Analyzers()
				expected := 5
				// Verifier le nombre d'analyseurs (001, 002, 003, 005, 006)
				if len(analyzers) != expected {
					t.Errorf("expected %d analyzers, got %d", expected, len(analyzers))
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, tt.testFunc)
	}
}
