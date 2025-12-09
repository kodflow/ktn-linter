// Good examples for the struct007 test case.
package struct007

// GoodService est une struct avec getter correctement nommé.
// Démontre la convention Go idiomatique: Name() retourne le champ name.
type GoodService struct {
	name string
}

// GoodServiceInterface définit le contrat public de GoodService.
type GoodServiceInterface interface {
	Name() string
}

// NewGoodService crée une nouvelle instance de GoodService.
//
// Params:
//   - name: nom du service
//
// Returns:
//   - *GoodService: nouvelle instance
func NewGoodService(name string) *GoodService {
	// Retour de la nouvelle instance
	return &GoodService{name: name}
}

// Name retourne le nom du service (getter correct: nom = champ).
//
// Returns:
//   - string: nom du service
func (s *GoodService) Name() string {
	// Retour du nom
	return s.name
}

// GoodDTO est un DTO sans getters - c'est acceptable pour les DTOs.
// Les DTOs n'ont pas besoin de getters car leurs champs sont publics.
type GoodDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// goodPrivateStruct est une struct privée (non exportée).
// La règle KTN-STRUCT-007 ne s'applique pas aux structs privées.
type goodPrivateStruct struct {
	value int
}

// GoodNoGetter est une struct sans getter pour son champ privé.
// C'est acceptable car les getters sont OPTIONNELS.
type GoodNoGetter struct {
	internalData string
}

// GoodNoGetterInterface définit les méthodes de GoodNoGetter.
type GoodNoGetterInterface interface {
	Process() error
}

// NewGoodNoGetter crée une nouvelle instance.
//
// Returns:
//   - *GoodNoGetter: nouvelle instance
func NewGoodNoGetter() *GoodNoGetter {
	// Retour de la nouvelle instance
	return &GoodNoGetter{}
}

// Process traite les données internes.
//
// Returns:
//   - error: erreur éventuelle
func (g *GoodNoGetter) Process() error {
	// Traitement
	return nil
}
