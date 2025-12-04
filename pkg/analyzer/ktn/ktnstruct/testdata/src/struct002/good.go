// Good examples for the struct002 test case.
package struct002

// UserService interface reprend toutes les méthodes publiques de userServiceImpl
type UserService interface {
	Create(name string) error
	GetByID(id int) (string, error)
	Update(id int, name string) error
}

// userServiceImpl implémente UserService avec toutes les méthodes
type userServiceImpl struct {
	users map[int]string
}

// useUserService utilise l'interface UserService.
//
// Params:
//   - us: service utilisateur
func useUserService(us UserService) {
	// Utilise l'interface
	_ = us
}

// Create crée un utilisateur.
//
// Params:
//   - name: nom de l'utilisateur
//
// Returns:
//   - error: erreur éventuelle
func (u *userServiceImpl) Create(name string) error {
	// Implementation
	return nil
}

// GetByID récupère un utilisateur par ID.
//
// Params:
//   - id: identifiant utilisateur
//
// Returns:
//   - string: nom de l'utilisateur
//   - error: erreur éventuelle
func (u *userServiceImpl) GetByID(id int) (string, error) {
	// Implementation
	return "", nil
}

// Update met à jour un utilisateur.
//
// Params:
//   - id: identifiant utilisateur
//   - name: nouveau nom
//
// Returns:
//   - error: erreur éventuelle
func (u *userServiceImpl) Update(id int, name string) error {
	// Implementation
	return nil
}

// helper est une fonction privée - pas besoin dans l'interface
func (u *userServiceImpl) helper() {
	// Private method
}

// Config est une struct simple sans méthode - PAS BESOIN D'INTERFACE.
// Contient la configuration de connexion réseau.
type Config struct {
	Host string
	Port int
}

// DataModel est une struct DTO sans méthode - PAS BESOIN D'INTERFACE.
// Représente un modèle de données avec identifiant, nom et tags.
type DataModel struct {
	ID   int
	Name string
	Tags []string
}

// init utilise les fonctions privées
func init() {
	// Appel de useUserService
	_ = useUserService(UserService{})
}
