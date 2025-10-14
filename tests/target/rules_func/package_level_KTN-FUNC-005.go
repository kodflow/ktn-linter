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
	Name    string
	Email   string
	Age     int
	Address string
	Phone   string
	Active  bool
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
		return 0, errors.New("name requis")
	}
	return 1, nil
}
