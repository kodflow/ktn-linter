package rules_func

import (
	"errors"
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-003: Commentaire godoc complet avec paramètres documentés
// ════════════════════════════════════════════════════════════════════════════

// CreateUserF003Good crée un nouvel utilisateur avec les informations fournies.
//
// Params:
//   - name: le nom de l'utilisateur
//   - email: l'adresse email de l'utilisateur
//   - age: l'âge de l'utilisateur
//
// Returns:
//   - int: l'ID du nouvel utilisateur créé
//   - error: une erreur si les données sont invalides
func CreateUserF003Good(name string, email string, age int) (int, error) {
	if name == "" {
		return 0, errors.New("name requis")
	}
	return 1, nil
}
