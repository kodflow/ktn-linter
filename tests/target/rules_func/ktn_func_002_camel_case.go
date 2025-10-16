package rules_func

import (
	"errors"
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-002: Fonction exportée avec commentaire godoc
// ════════════════════════════════════════════════════════════════════════════

// ProcessOrderF002Good traite une commande.
//
// Params:
//   - orderID: l'identifiant de la commande à traiter
//
// Returns:
//   - error: une erreur si l'orderID est invalide
func ProcessOrderF002Good(orderID int) error {
	if orderID <= 0 {
		// Retourne une erreur car l'orderID est invalide
		return errors.New("orderID invalide")
	}
	// Retourne nil car la commande est valide
	return nil
}

// ValidateEmailF002Good valide le format d'une adresse email.
//
// Params:
//   - email: l'adresse email à valider
//
// Returns:
//   - bool: true si l'email est valide, false sinon
func ValidateEmailF002Good(email string) bool {
	// Retourne true si l'email n'est pas vide
	return len(email) > 0
}
