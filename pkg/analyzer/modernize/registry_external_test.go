// External tests for registry.go - modernize package.
package modernize_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/modernize"
)

// TestAnalyzers tests the public Analyzers function.
func TestAnalyzers(t *testing.T) {
	tests := []struct {
		name      string
		wantEmpty bool
	}{
		{
			name:      "returns non-empty list",
			wantEmpty: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzers := modernize.Analyzers()
			// Vérification résultat
			if tt.wantEmpty && len(analyzers) > 0 {
				t.Error("Analyzers() returned non-empty when expected empty")
			}
			// Vérification non-vide
			if !tt.wantEmpty && len(analyzers) == 0 {
				t.Error("Analyzers() returned empty when expected non-empty")
			}
		})
	}
}

// TestAnalyzers_disabled tests that disabled analyzers are filtered out.
func TestAnalyzers_disabled(t *testing.T) {
	tests := []struct {
		name         string
		disabledName string
	}{
		{
			name:         "newexpr is disabled",
			disabledName: "newexpr",
		},
	}

	analyzers := modernize.Analyzers()
	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Vérifier que l'analyseur désactivé n'est pas présent
			for _, a := range analyzers {
				// Vérification de la condition
				if a.Name == tt.disabledName {
					t.Errorf("%s should be disabled but was found", tt.disabledName)
				}
			}
		})
	}
}
