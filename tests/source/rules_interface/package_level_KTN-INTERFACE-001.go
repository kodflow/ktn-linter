// Tests for KTN-INTERFACE-001: Package sans fichier interfaces.go
package rules_interface

// KTN-INTERFACE-001: Ce package devrait déclencher une erreur car il a une struct
// privée (implémentation) mais pas de fichier interfaces.go

// userServiceI001 est une implémentation privée qui devrait avoir une interface publique.
type userServiceI001 struct {
	db string
}

// GetUser récupère un utilisateur.
//
// Params:
//   - id: l'identifiant de l'utilisateur
//
// Returns:
//   - string: le nom de l'utilisateur
//   - error: une erreur si l'opération échoue
func (s *userServiceI001) GetUser(id string) (string, error) {
	return "user-" + id, nil
}

// NewUserServiceI001 crée une nouvelle instance du service.
//
// Params:
//   - db: la connexion à la base de données
//
// Returns:
//   - *userServiceI001: une nouvelle instance
func NewUserServiceI001(db string) *userServiceI001 {
	return &userServiceI001{db: db}
}
