// ════════════════════════════════════════════════════════════════════════════
// KTN-INTERFACE-001: Package sans fichier interfaces.go
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//
//	Un package contenant des structs avec méthodes (implémentations) doit avoir
//	un fichier interfaces.go définissant les interfaces publiques mockables.
//
//	POURQUOI :
//	- Facilite les tests unitaires avec des mocks
//	- Sépare contrats (interfaces) des implémentations
//	- Rend le code plus découplé et testable
//	- Convention Go pour architecture hexagonale/propre
//
// ❌ CAS INCORRECT (SEULE ERREUR: KTN-INTERFACE-001)
// NOTE: Ce package devrait avoir interfaces.go car il contient une struct
//
//	privée avec méthodes (implémentation mockable)
//
// ERREUR ATTENDUE: KTN-INTERFACE-001 sur le package
//
// ✅ CAS PARFAIT (voir target/) :
//
//	Créer interfaces.go avec:
//
//	// UserServiceI001Good définit le contrat du service utilisateur.
//	type UserServiceI001Good interface {
//	    GetUser(id string) (string, error)
//	}
//
//	func NewUserServiceI001Good(db string) UserServiceI001Good {
//	    return &userServiceImplI001Good{db: db}
//	}
//
// ════════════════════════════════════════════════════════════════════════════
package KTN_INTERFACE_001

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
