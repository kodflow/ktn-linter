// Tests for KTN-INTERFACE-001 (version corrigée)
package rules_interface

// KTN-INTERFACE-001 GOOD: Le package a maintenant interfaces.go

// userServiceImplI001Good est une implémentation privée avec interface publique.
type userServiceImplI001Good struct {
	db string
}

// GetUser implémente l'interface UserServiceI001Good.
//
// Params:
//   - id: l'identifiant de l'utilisateur
//
// Returns:
//   - string: le nom de l'utilisateur
//   - error: une erreur si l'opération échoue
func (s *userServiceImplI001Good) GetUser(id string) (string, error) {
	return "user-" + id, nil
}

// NewUserServiceI001Good crée une nouvelle instance du service.
//
// Params:
//   - db: la connexion à la base de données
//
// Returns:
//   - UserServiceI001Good: une nouvelle instance de l'interface
func NewUserServiceI001Good(db string) UserServiceI001Good {
	return &userServiceImplI001Good{db: db}
}
