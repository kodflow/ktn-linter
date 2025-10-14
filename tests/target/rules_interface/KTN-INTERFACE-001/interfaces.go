// ════════════════════════════════════════════════════════════════════════════
// KTN-INTERFACE-001: Package avec fichier interfaces.go (✅ CORRIGÉ)
// ════════════════════════════════════════════════════════════════════════════

package KTN_INTERFACE_001

// UserServiceI001Good définit le contrat du service utilisateur.
type UserServiceI001Good interface {
	// GetUser récupère un utilisateur.
	//
	// Params:
	//   - id: l'identifiant de l'utilisateur
	//
	// Returns:
	//   - string: le nom de l'utilisateur
	//   - error: une erreur si l'opération échoue
	GetUser(id string) (string, error)
}
