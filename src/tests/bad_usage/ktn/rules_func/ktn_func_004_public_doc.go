package rules_func

import (
	"context"
	"errors"
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-004: Commentaire godoc incomplet - retours non documentés
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1: Retours non documentés avec > 1 retour (SEULE ERREUR: KTN-FUNC-004)
// NOTE: Tout est parfait (nom + commentaire + params OK) SAUF retours non documentés
// ERREUR ATTENDUE: KTN-FUNC-004 sur FetchUserDataF004

// FetchUserDataF004 récupère les données utilisateur depuis une source externe.
//
// Params:
//   - ctx: le contexte pour gérer les timeouts et annulations
//   - userID: l'identifiant de l'utilisateur à récupérer
func FetchUserDataF004(ctx context.Context, userID int) (string, error) {
	if userID <= 0 {
		// Early return from function.
		return "", errors.New("userID invalide")
	}
	// Early return from function.
	return "user data", nil
}
