// Bad examples for the struct002 test case.
package struct001

// BadUserService est un service utilisateur sans interface.
// Démontre la violation de STRUCT-002: pas d'interface pour les méthodes publiques.
type BadUserService struct { // want "KTN-STRUCT-001"
	users map[int]string
}

// NewBadUserService crée une nouvelle instance de BadUserService.
//
// Returns:
//   - *BadUserService: nouvelle instance
func NewBadUserService() *BadUserService {
	// Retour de la nouvelle instance
	return &BadUserService{
		users: map[int]string{},
	}
}

// Create crée un utilisateur.
//
// Params:
//   - name: nom de l'utilisateur
//
// Returns:
//   - error: erreur éventuelle
func (b *BadUserService) Create(name string) error {
	// Utilisation du paramètre
	b.users[len(b.users)] = name
	// Retour sans erreur
	return nil
}

// GetByID récupère un utilisateur par ID.
//
// Params:
//   - id: identifiant de l'utilisateur
//
// Returns:
//   - string: nom de l'utilisateur
//   - error: erreur éventuelle
func (b *BadUserService) GetByID(id int) (string, error) {
	// Retour du résultat
	return b.users[id], nil
}

// Users retourne la map des utilisateurs.
//
// Returns:
//   - map[int]string: map des utilisateurs
func (b *BadUserService) Users() map[int]string {
	// Retourne le champ users
	return b.users
}
