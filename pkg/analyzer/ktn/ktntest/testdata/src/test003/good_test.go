package test003_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest/testdata/src/test002"
)

// TestProcessItemComplete teste avec couverture complète (BIEN).
// A un fichier source good.go correspondant.
func TestProcessItemComplete(t *testing.T) {
	// Table de tests
	tests := []struct {
		name      string
		item      string
		want      string
		wantError bool
	}{
		{"valid item", "data", "processed:data", false},
		{"empty item", "", "", true},
	}

	// Parcours des tests
	for _, tt := range tests {
		// Exécution sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Appel fonction
			result, err := test002.ProcessItem(tt.item)
			// Vérification erreur
			if (err != nil) != tt.wantError {
				t.Errorf("ProcessItem(%q) error = %v, wantError %v", tt.item, err, tt.wantError)
			}
			// Vérification résultat
			if !tt.wantError && result != tt.want {
				t.Errorf("ProcessItem(%q) = %q, want %q", tt.item, result, tt.want)
			}
		})
	}
}

// TestCountItemsComplete teste avec couverture complète (BIEN).
// A un fichier source good.go correspondant.
func TestCountItemsComplete(t *testing.T) {
	const MAX_ITEMS int = 1000
	// Table de tests
	tests := []struct {
		name      string
		items     []string
		want      int
		wantError bool
	}{
		{"empty list", []string{}, 0, false},
		{"single item", []string{"a"}, 1, false},
		{"multiple items", []string{"a", "b", "c"}, 3, false},
		{"too many items", make([]string, MAX_ITEMS+1), 0, true},
	}

	// Parcours des tests
	for _, tt := range tests {
		// Exécution sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Appel fonction
			result, err := test002.CountItems(tt.items)
			// Vérification erreur
			if (err != nil) != tt.wantError {
				t.Errorf("CountItems() error = %v, wantError %v", err, tt.wantError)
			}
			// Vérification résultat
			if !tt.wantError && result != tt.want {
				t.Errorf("CountItems() = %d, want %d", result, tt.want)
			}
		})
	}
}
