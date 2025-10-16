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
		// Retourne 0 et une erreur car le nom est requis
		return 0, errors.New("name requis")
	}
	// Retourne l'ID de l'utilisateur créé et nil pour l'erreur
	return 1, nil
}
