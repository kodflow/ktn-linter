// Tests for KTN-INTERFACE-002: Type public défini comme struct au lieu d'interface
package rules_interface

// KTN-INTERFACE-002: UserServiceI002 est public mais c'est une struct
// Il devrait être une interface dans interfaces.go

// UserServiceI002 devrait être une interface, pas une struct publique.
type UserServiceI002 struct {
	DB       string
	CacheDir string
}

// ProcessI002 traite des données.
//
// Params:
//   - data: les données à traiter
//
// Returns:
//   - error: une erreur si l'opération échoue
func (s *UserServiceI002) ProcessI002(data string) error {
	return nil
}

// OrderManagerI002 est aussi une struct publique (violation).
type OrderManagerI002 struct {
	Orders []string
}

// CreateOrderI002 crée une commande.
//
// Params:
//   - order: la commande à créer
//
// Returns:
//   - error: une erreur si l'opération échoue
func (o *OrderManagerI002) CreateOrderI002(order string) error {
	o.Orders = append(o.Orders, order)
	return nil
}
