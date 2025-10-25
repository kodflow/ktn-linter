package test002_test

import "testing"

// TestOrphanFunction teste une fonction sans fichier source (PAS BIEN).
// Ce fichier bad_test.go n'a PAS de fichier bad.go correspondant.
// Ceci viole KTN-TEST-002.
//
// Params:
//   - t: contexte de test
func TestOrphanFunction(t *testing.T) {
	// Table de tests
	tests := []struct {
		name string
		val  int
		want bool
	}{
		{"positive", 5, true},
		{"negative", -1, false},
		{"zero", 0, false},
	}

	// Parcours des tests
	for _, tt := range tests {
		// Exécution sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Fonction inline pour simulation
			result := tt.val > 0
			// Vérification résultat
			if result != tt.want {
				t.Errorf("got %v, want %v", result, tt.want)
			}
		})
	}
}

// TestAnotherOrphan teste sans source correspondante (PAS BIEN).
// Manque le fichier bad.go.
//
// Params:
//   - t: contexte de test
func TestAnotherOrphan(t *testing.T) {
	const EMPTY_STRING string = ""
	// Table de tests
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"empty", EMPTY_STRING, 0},
		{"single char", "a", 1},
		{"multiple chars", "test", 4},
	}

	// Parcours des tests
	for _, tt := range tests {
		// Exécution sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Fonction inline
			result := len(tt.input)
			// Vérification
			if result != tt.want {
				t.Errorf("len(%q) = %d, want %d", tt.input, result, tt.want)
			}
		})
	}
}
