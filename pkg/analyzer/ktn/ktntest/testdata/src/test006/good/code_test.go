package test006

import "testing"

// TestProcessDataWithErrors teste avec couverture complète (BIEN).
// Fichier code.go existe dans le même package - pattern 1:1 respecté.
//
// Params:
//   - t: contexte de test
func TestProcessDataWithErrors(t *testing.T) {
	const MAX_DATA_LENGTH int = 100
	// Table de tests
	tests := []struct {
		name      string
		data      string
		want      string
		wantError bool
	}{
		{"valid data", "test", "processed:test", false},
		{"empty data", "", "processed:", false},
		{"too long data", string(make([]byte, MAX_DATA_LENGTH+1)), "", true},
	}

	// Parcours des tests
	for _, tt := range tests {
		// Exécution sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Appel fonction privée
			result, err := processData(tt.data)
			// Vérification erreur
			if (err != nil) != tt.wantError {
				t.Errorf("processData() error = %v, wantError %v", err, tt.wantError)
			}
			// Vérification résultat
			if !tt.wantError && result != tt.want {
				t.Errorf("processData() = %q, want %q", result, tt.want)
			}
		})
	}
}

// TestValidateNumberWithErrors teste avec couverture complète (BIEN).
// Fichier code.go existe dans le même package - pattern 1:1 respecté.
//
// Params:
//   - t: contexte de test
func TestValidateNumberWithErrors(t *testing.T) {
	// Table de tests
	tests := []struct {
		name      string
		num       int
		wantError bool
	}{
		{"positive number", 5, false},
		{"zero", 0, false},
		{"negative number", -1, true},
	}

	// Parcours des tests
	for _, tt := range tests {
		// Exécution sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Appel fonction privée
			err := validateNumber(tt.num)
			// Vérification erreur
			if (err != nil) != tt.wantError {
				t.Errorf("validateNumber(%d) error = %v, wantError %v", tt.num, err, tt.wantError)
			}
		})
	}
}
