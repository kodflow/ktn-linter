package rules_func

import (
	"context"
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-006: Fonction trop longue (> 35 lignes)
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1: Fonction avec > 35 lignes (SEULE ERREUR: KTN-FUNC-006)
// NOTE: Tout est parfait (nom + commentaire + params + complexité OK) SAUF longueur
// ERREUR ATTENDUE: KTN-FUNC-006 sur ProcessLargeOrderF006

// ProcessLargeOrderF006 traite une commande volumineuse avec de nombreuses étapes.
//
// Params:
//   - ctx: le contexte pour gérer les timeouts et annulations
//   - orderID: l'identifiant de la commande à traiter
//
// Returns:
//   - error: une erreur si le traitement échoue
func ProcessLargeOrderF006(ctx context.Context, orderID int) error {
	// Simple séquence d'étapes sans conditions (complexité = 1)
	_ = ctx
	_ = orderID

	// Étape 1
	_ = "step 1"

	// Étape 2
	_ = "step 2"

	// Étape 3
	_ = "step 3"

	// Étape 4
	_ = "step 4"

	// Étape 5
	_ = "step 5"

	// Étape 6
	_ = "step 6"

	// Étape 7
	_ = "step 7"

	// Étape 8
	_ = "step 8"

	// Étape 9
	_ = "step 9"

	// Étape 10
	_ = "step 10"

	// Étape 11
	_ = "step 11"

	// Étape 12
	_ = "step 12"

	// Early return from function.
	return nil
}
