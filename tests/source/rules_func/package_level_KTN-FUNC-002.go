package rules_func

import (
	"errors"
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-002: Fonction exportée sans commentaire godoc
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1: Fonction exportée sans commentaire (SEULE ERREUR: KTN-FUNC-002)
// NOTE: Nom OK, params OK, longueur OK, MAIS pas de commentaire godoc
// ERREUR ATTENDUE: KTN-FUNC-002 sur ProcessOrderF002

func ProcessOrderF002(orderID int) error {
	if orderID <= 0 {
		return errors.New("orderID invalide")
	}
	return nil
}

// ❌ CAS INCORRECT 2: Autre fonction exportée sans commentaire (SEULE ERREUR: KTN-FUNC-002)
// ERREUR ATTENDUE: KTN-FUNC-002 sur ValidateEmailF002

func ValidateEmailF002(email string) bool {
	return len(email) > 0
}
