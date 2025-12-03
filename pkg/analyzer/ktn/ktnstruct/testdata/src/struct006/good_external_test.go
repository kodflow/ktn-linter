// Package struct006_test contient les tests black-box pour struct006.
package struct006_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct/testdata/src/struct006"
)

// TestGoodUser vérifie le fonctionnement de GoodUser.
func TestGoodUser(t *testing.T) {
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
			user := struct006.NewGoodUser(tt.id, tt.userName, tt.email)

			// Vérification ID
			if user.ID() != tt.expectedID {
				// Erreur si ID incorrect
				t.Errorf("ID() = %d, want %d", user.ID(), tt.expectedID)
			}

			// Vérification Name
			if user.Name() != tt.expectedName {
				// Erreur si nom incorrect
				t.Errorf("Name() = %s, want %s", user.Name(), tt.expectedName)
			}

			// Vérification Email
			if user.Email() != tt.expectedEmail {
				// Erreur si email incorrect
				t.Errorf("Email() = %s, want %s", user.Email(), tt.expectedEmail)
			}
		})
	}
}

// TestGoodUser_SetName vérifie le setter de nom.
func TestGoodUser_SetName(t *testing.T) {
	// Table des cas de test
	tests := []struct {
		name        string
		initialName string
		newName     string
	}{
		{
			name:        "changement de nom simple",
			initialName: "John",
			newName:     "Jane",
		},
		{
			name:        "changement vers nom vide",
			initialName: "John",
			newName:     "",
		},
	}

	// Exécution des tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Création de l'utilisateur
			user := struct006.NewGoodUser(1, tt.initialName, "test@example.com")

			// Modification du nom
			user.SetName(tt.newName)

			// Vérification du nouveau nom
			if user.Name() != tt.newName {
				// Erreur si nom incorrect
				t.Errorf("Name() = %s, want %s", user.Name(), tt.newName)
			}
		})
	}
}
