// ════════════════════════════════════════════════════════════════════════════
// KTN-INTERFACE-002: Type public défini comme struct au lieu d'interface
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//
//	Les types publics doivent être des interfaces (dans interfaces.go),
//	pas des structs. Les structs doivent être privées (implémentations).
//
//	POURQUOI :
//	- Permet le mocking facile dans les tests
//	- Découple les dépendances entre packages
//	- Force une architecture orientée contrats
//	- Facilite les changements d'implémentation
//
// ❌ CAS INCORRECT 1: Struct publique (SEULE ERREUR: KTN-INTERFACE-002)
// NOTE: UserServiceI002 devrait être une interface dans interfaces.go
// ERREUR ATTENDUE: KTN-INTERFACE-002 sur UserServiceI002
//
// ❌ CAS INCORRECT 2: Struct publique (SEULE ERREUR: KTN-INTERFACE-002)
// NOTE: OrderManagerI002 devrait être une interface dans interfaces.go
// ERREUR ATTENDUE: KTN-INTERFACE-002 sur OrderManagerI002
//
// ✅ CAS PARFAIT (voir target/) :
//
//	// interfaces.go
//	type UserService interface {
//	    Process(data string) error
//	}
//
//	// impl.go (implémentation privée)
//	type userServiceImpl struct { ... }
//
// ════════════════════════════════════════════════════════════════════════════
package KTN_INTERFACE_002

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
