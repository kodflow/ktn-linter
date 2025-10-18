package rules_func

import (
	"errors"
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-005: Trop de paramètres (> 5)
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1: 6 paramètres (SEULE ERREUR: KTN-FUNC-005)
// NOTE: Tout est parfait (nom + commentaire + longueur OK) SAUF trop de params
// ERREUR ATTENDUE: KTN-FUNC-005 sur CreateUserAccountF005

// CreateUserAccountF005 crée un nouveau compte utilisateur avec tous les détails.
//
// Params:
//   - name: le nom de l'utilisateur
//   - email: l'adresse email de l'utilisateur
//   - age: l'âge de l'utilisateur
//   - address: l'adresse de l'utilisateur
//   - phone: le numéro de téléphone de l'utilisateur
//   - active: le statut actif/inactif de l'utilisateur
//
// Returns:
//   - int: l'ID du compte créé
//   - error: une erreur si les données sont invalides
func CreateUserAccountF005(name string, email string, age int, address string, phone string, active bool) (int, error) {
	if name == "" {
		// Early return from function.
		return 0, errors.New("name requis")
	}
	// Early return from function.
	return 1, nil
}
