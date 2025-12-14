package ktnstruct_test

import (
	"testing"

	ktnstruct "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
)

// TestStruct001Deprecated vérifie que l'analyzer déprécié existe toujours pour compatibilité.
// DEPRECATED: KTN-STRUCT-001 est déprécié et remplacé par KTN-API-001.
// Les "mirror interfaces" (100% des méthodes) sont un anti-pattern.
// Utilisez KTN-API-001 pour les interfaces minimales côté consumer (ISP).
func TestStruct001Deprecated(t *testing.T) {
	tests := []struct {
		name           string
		checkDeprecated bool
	}{
		{
			name:           "analyzer_exists_for_backward_compatibility",
			checkDeprecated: true,
		},
	}

	// Parcourir les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Vérifier que l'analyzer existe toujours pour compatibilité
			if tt.checkDeprecated && ktnstruct.Analyzer001 == nil {
				t.Error("Analyzer001 should exist for backward compatibility")
			}
			// Note: l'analyzer n'est plus dans le registry (GetAnalyzers)
			// mais le code existe pour documentation et migration
		})
	}
}
