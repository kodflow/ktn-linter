// ════════════════════════════════════════════════════════════════════════════
// KTN-INTERFACE-002: Types publics sont des interfaces (✅ CORRIGÉ)
// ════════════════════════════════════════════════════════════════════════════

package KTN_INTERFACE_002

// UserServiceI002Good définit le contrat du service utilisateur.
type UserServiceI002Good interface {
	// ProcessI002Good traite des données.
	//
	// Params:
	//   - data: les données à traiter
	//
	// Returns:
	//   - error: une erreur si l'opération échoue
	ProcessI002Good(data string) error
}

// OrderManagerI002Good définit le contrat de gestion des commandes.
type OrderManagerI002Good interface {
	// CreateOrderI002Good crée une commande.
	//
	// Params:
	//   - order: la commande à créer
	//
	// Returns:
	//   - error: une erreur si l'opération échoue
	CreateOrderI002Good(order string) error
}
