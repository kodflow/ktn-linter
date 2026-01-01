// Package struct001 provides good test cases.
// GoodService demonstrates correct getter naming.
package struct001

// GoodService est une struct avec getter correctement nommé.
// Démontre la convention Go idiomatique: Name() retourne le champ name.
type GoodService struct {
	name string
	age  int
}

// GoodServiceInterface définit le contrat public de GoodService.
type GoodServiceInterface interface {
	Name() string
	Age() int
}

// NewGoodService crée une nouvelle instance de GoodService.
//
// Params:
//   - name: nom du service
//   - age: âge du service
//
// Returns:
//   - *GoodService: nouvelle instance
func NewGoodService(name string, age int) *GoodService {
	// Retour de la nouvelle instance
	return &GoodService{name: name, age: age}
}

// Name retourne le nom du service (getter correct: Name = champ name).
//
// Returns:
//   - string: nom du service
func (s *GoodService) Name() string {
	// Retour du nom
	return s.name
}

// Age retourne l'âge du service (getter correct: Age = champ age).
//
// Returns:
//   - int: âge du service
func (s *GoodService) Age() int {
	// Retour de l'âge
	return s.age
}
