package test006

import "testing"

// TestOrphanFunction teste sans fichier source (PAS BIEN).
// Pas de fichier orphan.go - viole KTN-TEST-006.
func TestOrphanFunction(t *testing.T) {
	// Table de tests
	tests := []struct {
		name  string
		value int
		want  bool
	}{
		{"positive", 10, true},
		{"zero", 0, false},
		{"negative", -5, false},
	}

	// Parcours des tests
	for _, tt := range tests {
		// Exécution sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Fonction inline (pas de fichier source)
			result := tt.value > 0
			// Vérification
			if result != tt.want {
				t.Errorf("got %v, want %v", result, tt.want)
			}
		})
	}
}

// TestAnotherOrphan teste sans source (PAS BIEN).
// Manque fichier orphan.go correspondant.
func TestAnotherOrphan(t *testing.T) {
	const EMPTY_INPUT string = ""
	// Table de tests
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"empty", EMPTY_INPUT, 0},
		{"single", "a", 1},
		{"multiple", "hello", 5},
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
