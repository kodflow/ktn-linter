// External tests for ktnapi registry.
package ktnapi_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnapi"
)

// TestAnalyzers tests that Analyzers returns all API analyzers.
func TestAnalyzers(t *testing.T) {
	tests := []struct {
		name          string
		expectedCount int
		expectedName  string
	}{
		{
			name:          "returns_exactly_one_analyzer",
			expectedCount: 1,
			expectedName:  "ktnapi001",
		},
	}

	// Parcourir les tests
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			analyzers := ktnapi.Analyzers()
			// Vérification du nombre d'analyzers
			if len(analyzers) != tt.expectedCount {
				t.Errorf("expected %d analyzer, got %d", tt.expectedCount, len(analyzers))
				return
			}

			// Vérification de l'analyzer KTN-API-001
			if analyzers[0].Name != tt.expectedName {
				t.Errorf("expected %s, got %s", tt.expectedName, analyzers[0].Name)
			}
		})
	}
}
