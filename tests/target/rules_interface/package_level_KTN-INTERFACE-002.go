// Tests for KTN-INTERFACE-002 (version corrigée)
package rules_interface

// KTN-INTERFACE-002 GOOD: Les types publics sont maintenant des interfaces

// userServiceImplI002Good est l'implémentation privée.
type userServiceImplI002Good struct {
	db       string
	cacheDir string
}

// ProcessI002Good traite des données.
//
// Params:
//   - data: les données à traiter
//
// Returns:
//   - error: une erreur si l'opération échoue
func (s *userServiceImplI002Good) ProcessI002Good(data string) error {
	return nil
}

// orderManagerImplI002Good est l'implémentation privée.
type orderManagerImplI002Good struct {
	orders []string
}

// CreateOrderI002Good crée une commande.
//
// Params:
//   - order: la commande à créer
//
// Returns:
//   - error: une erreur si l'opération échoue
func (o *orderManagerImplI002Good) CreateOrderI002Good(order string) error {
	o.orders = append(o.orders, order)
	return nil
}
