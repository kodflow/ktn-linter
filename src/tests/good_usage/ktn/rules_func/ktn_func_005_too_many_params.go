package rules_func

import (
	"context"
	"errors"
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-005: Fonction avec ≤ 5 paramètres
// ════════════════════════════════════════════════════════════════════════════

// UserF005GoodConfig contient les détails d'un utilisateur (✅ struct pour éviter trop de params)
type UserF005GoodConfig struct {
	// Name est le nom de l'utilisateur
	Name string
	// Email est l'adresse email de l'utilisateur
	Email string
	// Age est l'âge de l'utilisateur
	Age int
	// Address est l'adresse postale de l'utilisateur
	Address string
	// Phone est le numéro de téléphone de l'utilisateur
	Phone string
	// Active indique si le compte utilisateur est actif
	Active bool
}

// CreateUserAccountF005Good crée un nouveau compte utilisateur avec la configuration fournie.
//
// Params:
//   - ctx: le contexte pour gérer les timeouts
//   - config: la configuration de l'utilisateur
//
// Returns:
//   - int: l'ID du compte créé
//   - error: une erreur si les données sont invalides
func CreateUserAccountF005Good(ctx context.Context, config UserF005GoodConfig) (int, error) {
	if config.Name == "" {
		// Retourne 0 et une erreur car le nom est requis
		return 0, errors.New("name requis")
	}
	// Retourne l'ID du compte créé et nil pour l'erreur
	return 1, nil
}
