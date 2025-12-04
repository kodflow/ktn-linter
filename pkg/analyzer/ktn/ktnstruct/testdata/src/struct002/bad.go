// Bad examples for the struct005 test case.
package struct002

// BadUserService est un service utilisateur sans constructeur.
// Démontre la violation STRUCT-005: pas de NewBadUserService().
type BadUserService struct { // want "KTN-STRUCT-002"
	users map[int]string
}

// BadUserServiceInterface définit les méthodes de BadUserService.
type BadUserServiceInterface interface {
	Create(name string) error
	GetByID(id int) string
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
func (b *BadUserService) GetByID(id int) string {
	// Retour du résultat
	return b.users[id]
}
