// ════════════════════════════════════════════════════════════════════════════
// KTN-TEST-001: Package name incorrect dans fichier *_test.go
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//
//	Les fichiers *_test.go doivent utiliser le package name avec suffixe _test
//	(ex: mypackage_test) pour forcer les tests à utiliser l'API publique.
//
//	POURQUOI :
//	- Force les tests à utiliser uniquement l'API publique
//	- Empêche les couplages trop forts avec l'implémentation interne
//	- Convention Go standard (black-box testing)
//	- Détecte les problèmes d'encapsulation
//
// ❌ CAS INCORRECT : Voir service_test.go avec mauvais package
// ERREUR ATTENDUE: KTN-TEST-001 sur service_test.go
//
// ✅ CAS PARFAIT (voir target/) :
//
//	// service_test.go
//	package KTN_TEST_001_test  // ✓ Avec _test
//
// ════════════════════════════════════════════════════════════════════════════
package KTN_TEST_001

// UserService gère les utilisateurs.
type UserService struct {
	users map[string]string
}

// NewUserService crée un nouveau service.
//
// Returns:
//   - *UserService: nouvelle instance
func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]string),
	}
}

// AddUser ajoute un utilisateur.
//
// Params:
//   - id: identifiant de l'utilisateur
//   - name: nom de l'utilisateur
func (s *UserService) AddUser(id string, name string) {
	s.users[id] = name
}

// GetUser récupère un utilisateur.
//
// Params:
//   - id: identifiant de l'utilisateur
//
// Returns:
//   - string: nom de l'utilisateur
//   - bool: true si trouvé
func (s *UserService) GetUser(id string) (string, bool) {
	name, exists := s.users[id]
	return name, exists
}
