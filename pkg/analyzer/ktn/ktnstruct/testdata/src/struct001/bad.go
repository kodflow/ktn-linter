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

// BadPartialInterface est une interface incomplète (ne couvre pas toutes les méthodes).
type BadPartialInterface interface {
	Process() error
}

// BadIncompleteImpl a un compile-time check mais l'interface ne couvre pas toutes les méthodes.
// Cela doit TOUJOURS déclencher KTN-STRUCT-001 car l'interface est incomplète.
type BadIncompleteImpl struct { // want "KTN-STRUCT-001"
	data string
}

// Process implémente BadPartialInterface.
//
// Returns:
//   - error: erreur éventuelle
func (b *BadIncompleteImpl) Process() error {
	// Implementation
	return nil
}

// ExtraMethod est une méthode publique NON couverte par BadPartialInterface.
// Cela prouve que le compile-time check ne suffit pas - l'interface doit être complète.
//
// Returns:
//   - string: données
func (b *BadIncompleteImpl) ExtraMethod() string {
	// Retour des données
	return b.data
}

// Compile-time check - MAIS l'interface ne couvre pas ExtraMethod()
var _ BadPartialInterface = (*BadIncompleteImpl)(nil)
