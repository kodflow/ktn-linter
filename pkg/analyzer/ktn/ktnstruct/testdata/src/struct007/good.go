// Good examples for the struct007 test case.
package struct007

// GoodService est une struct non-DTO avec getter pour son champ privé.
// Démonstration d'un service avec encapsulation correcte.
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

// Name retourne le nom du service.
//
// Returns:
//   - string: nom du service
func (s *GoodService) Name() string {
	// Retour du nom
	return s.name
}

// GoodDTO est un DTO qui n'a pas besoin de getters.
// Les DTOs sont exemptés de la règle KTN-STRUCT-007.
type GoodDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// goodPrivateStruct est une struct privée (non exportée).
// La règle KTN-STRUCT-007 ne s'applique pas aux structs privées.
type goodPrivateStruct struct {
	value int
}

// GoodAllPublic est une struct avec tous champs publics.
// Pas de champs privés = pas besoin de getters.
type GoodAllPublic struct {
	Name  string
	Value int
}
