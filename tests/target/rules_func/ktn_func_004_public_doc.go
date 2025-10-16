package rules_func

import (
	"context"
	"errors"
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-004: Commentaire godoc complet avec retours documentés
// ════════════════════════════════════════════════════════════════════════════

// FetchUserDataF004Good récupère les données utilisateur depuis une source externe.
//
// Params:
//   - ctx: le contexte pour gérer les timeouts et annulations
//   - userID: l'identifiant de l'utilisateur à récupérer
//
// Returns:
//   - string: les données utilisateur récupérées
//   - error: une erreur si l'utilisateur n'existe pas ou si le contexte expire
func FetchUserDataF004Good(ctx context.Context, userID int) (string, error) {
	if userID <= 0 {
		// Retourne une chaîne vide et une erreur car l'userID est invalide
		return "", errors.New("userID invalide")
	}
	// Retourne les données utilisateur et nil pour l'erreur
	return "user data", nil
}
