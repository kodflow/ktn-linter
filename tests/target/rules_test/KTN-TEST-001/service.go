// ════════════════════════════════════════════════════════════════════════════
// KTN-TEST-001: Package name correct avec suffixe _test (✅ CORRIGÉ)
// ════════════════════════════════════════════════════════════════════════════
package KTN_TEST_001_GOOD

// UserServiceData gère les utilisateurs (Data suffix autorisé pour struct publique).
type UserServiceData struct {
	users map[string]string
}

// NewUserServiceData crée un nouveau service.
//
// Returns:
//   - *UserServiceData: nouvelle instance
func NewUserServiceData() *UserServiceData {
	return &UserServiceData{
		users: make(map[string]string),
	}
}

// AddUser ajoute un utilisateur.
//
// Params:
//   - id: identifiant de l'utilisateur
//   - name: nom de l'utilisateur
func (s *UserServiceData) AddUser(id string, name string) {
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
func (s *UserServiceData) GetUser(id string) (string, bool) {
	name, exists := s.users[id]
	return name, exists
}
