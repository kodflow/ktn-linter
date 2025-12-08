// Bad examples for the struct007 test case.
package struct007

// BadService est une struct non-DTO sans getter pour son champ privé.
// Démontre la violation STRUCT-007: pas de getter Name() pour le champ name.
type BadService struct {
	name string // want "KTN-STRUCT-007: la struct 'BadService' devrait avoir un getter 'Name\\(\\)' pour le champ privé 'name'"
}

// BadServiceInterface définit les méthodes de BadService.
type BadServiceInterface interface {
	Run() error
}

// NewBadService crée une nouvelle instance de BadService.
//
// Returns:
//   - *BadService: nouvelle instance
func NewBadService() *BadService {
	// Retourne une nouvelle instance
	return &BadService{}
}

// Run exécute le service.
//
// Returns:
//   - error: erreur éventuelle
func (s *BadService) Run() error {
	// Retourne nil si succès
	return nil
}

// BadRepository est une struct non-DTO avec plusieurs champs privés sans getters.
// Démontre la violation STRUCT-007 avec plusieurs champs privés.
type BadRepository struct {
	connection string // want "KTN-STRUCT-007: la struct 'BadRepository' devrait avoir un getter 'Connection\\(\\)' pour le champ privé 'connection'"
	timeout    int    // want "KTN-STRUCT-007: la struct 'BadRepository' devrait avoir un getter 'Timeout\\(\\)' pour le champ privé 'timeout'"
}

// BadRepositoryInterface définit les méthodes de BadRepository.
type BadRepositoryInterface interface {
	Connect() error
}

// NewBadRepository crée une nouvelle instance de BadRepository.
//
// Returns:
//   - *BadRepository: nouvelle instance
func NewBadRepository() *BadRepository {
	// Retourne une nouvelle instance
	return &BadRepository{}
}

// Connect établit la connexion.
//
// Returns:
//   - error: erreur éventuelle
func (r *BadRepository) Connect() error {
	// Retourne nil si succès
	return nil
}

// BadHandler est une struct avec un champ privé mais pas le bon getter.
// Démontre la violation STRUCT-007 avec un getter mal nommé.
type BadHandler struct {
	logger string // want "KTN-STRUCT-007: la struct 'BadHandler' devrait avoir un getter 'Logger\\(\\)' pour le champ privé 'logger'"
}

// BadHandlerInterface définit les méthodes de BadHandler.
type BadHandlerInterface interface {
	WrongGetter() int
}

// NewBadHandler crée une nouvelle instance de BadHandler.
//
// Returns:
//   - *BadHandler: nouvelle instance
func NewBadHandler() *BadHandler {
	// Retourne une nouvelle instance
	return &BadHandler{}
}

// WrongGetter retourne autre chose - pas le bon getter.
//
// Returns:
//   - int: toujours 0
func (h *BadHandler) WrongGetter() int {
	// Retour de la valeur
	return 0
}
