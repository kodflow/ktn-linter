package test002_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest/testdata/src/test001"
)

// TestValidateWithErrors teste avec couverture complète (BIEN).
// Utilise le package externe test001_test et table-driven tests.
//
// Params:
//   - t: contexte de test
func TestValidateWithErrors(t *testing.T) {
	// Table de tests
	tests := []struct {
		name      string
		value     int
		wantError bool
	}{
		{"valid value", 50, false},
		{"negative value", -1, true},
		{"exceeds maximum", 150, true},
	}

	// Parcours des tests
	for _, tt := range tests {
		// Exécution sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Appel fonction
			err := test001.Validate(tt.value)
			// Vérification erreur
			if (err != nil) != tt.wantError {
				t.Errorf("Validate(%d) error = %v, wantError %v", tt.value, err, tt.wantError)
			}
		})
	}
}

// TestTransformWithErrors teste avec couverture complète (BIEN).
// Utilise le package externe test001_test et table-driven tests.
//
// Params:
//   - t: contexte de test
func TestTransformWithErrors(t *testing.T) {
	// Table de tests
	tests := []struct {
		name      string
		input     string
		want      string
		wantError bool
	}{
		{"valid input", "test", "[test]", false},
		{"empty string", "", "", true},
	}

	// Parcours des tests
	for _, tt := range tests {
		// Exécution sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Appel fonction
			result, err := test001.Transform(tt.input)
			// Vérification erreur
			if (err != nil) != tt.wantError {
				t.Errorf("Transform(%q) error = %v, wantError %v", tt.input, err, tt.wantError)
			}
			// Vérification résultat si pas d'erreur
			if !tt.wantError && result != tt.want {
				t.Errorf("Transform(%q) = %q, want %q", tt.input, result, tt.want)
			}
		})
	}
}
