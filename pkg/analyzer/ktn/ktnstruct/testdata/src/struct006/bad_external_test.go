// Package struct006_test contient les tests black-box pour struct006.
package struct006_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct/testdata/src/struct006"
)

// TestBadUser vérifie le fonctionnement de BadUser.
func TestBadUser(t *testing.T) {
	// Table des cas de test
	tests := []struct {
		name          string
		id            int
		userName      string
		email         string
		expectedID    int
		expectedName  string
		expectedEmail string
	}{
		{
			name:          "création utilisateur standard",
			id:            1,
			userName:      "John",
			email:         "john@example.com",
			expectedID:    1,
			expectedName:  "John",
			expectedEmail: "john@example.com",
		},
		{
			name:          "création utilisateur avec id différent",
			id:            42,
			userName:      "Jane",
			email:         "jane@example.com",
			expectedID:    42,
			expectedName:  "Jane",
			expectedEmail: "jane@example.com",
		},
	}

	// Exécution des tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Création de l'utilisateur
			user := struct006.NewBadUser(tt.id, tt.userName, tt.email)

			// Vérification ID (via getter non idiomatique)
			if user.GetID() != tt.expectedID {
				// Erreur si ID incorrect
				t.Errorf("GetID() = %d, want %d", user.GetID(), tt.expectedID)
			}

			// Vérification Name (via getter non idiomatique)
			if user.GetName() != tt.expectedName {
				// Erreur si nom incorrect
				t.Errorf("GetName() = %s, want %s", user.GetName(), tt.expectedName)
			}

			// Vérification Email (via getter non idiomatique)
			if user.GetEmail() != tt.expectedEmail {
				// Erreur si email incorrect
				t.Errorf("GetEmail() = %s, want %s", user.GetEmail(), tt.expectedEmail)
			}
		})
	}
}

// TestBadUser_Save vérifie la sauvegarde.
func TestBadUser_Save(t *testing.T) {
	// Table des cas de test
	tests := []struct {
		name      string
		id        int
		userName  string
		email     string
		expectErr bool
	}{
		{
			name:      "sauvegarde utilisateur standard",
			id:        1,
			userName:  "John",
			email:     "john@example.com",
			expectErr: false,
		},
		{
			name:      "sauvegarde autre utilisateur",
			id:        2,
			userName:  "Jane",
			email:     "jane@example.com",
			expectErr: false,
		},
	}

	// Exécution des tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Création de l'utilisateur
			user := struct006.NewBadUser(tt.id, tt.userName, tt.email)

			// Sauvegarde
			err := user.Save()

			// Vérification de l'erreur
			if (err != nil) != tt.expectErr {
				// Erreur si résultat inattendu
				t.Errorf("Save() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}
