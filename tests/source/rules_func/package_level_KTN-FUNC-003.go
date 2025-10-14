package rules_func

import (
	"errors"
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-003: Commentaire godoc incomplet - paramètres non documentés
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1: Params non documentés avec > 2 params (SEULE ERREUR: KTN-FUNC-003)
// NOTE: Tout est parfait (nom + commentaire présent + params OK) SAUF params non mentionnés dans doc
// ERREUR ATTENDUE: KTN-FUNC-003 sur CreateUserF003

// CreateUserF003 crée un nouvel utilisateur.
//
// Returns:
//   - int: l'ID du nouvel utilisateur créé
//   - error: une erreur si les données sont invalides
func CreateUserF003(name string, email string, age int) (int, error) {
	if name == "" {
		return 0, errors.New("name requis")
	}
	return 1, nil
}
