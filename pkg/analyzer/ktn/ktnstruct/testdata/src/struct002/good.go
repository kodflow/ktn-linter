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

// Create crée un utilisateur
func (u *userServiceImpl) Create(name string) error {
	// Implementation
	return nil
}

// GetByID récupère un utilisateur par ID
func (u *userServiceImpl) GetByID(id int) (string, error) {
	// Implementation
	return "", nil
}

// Update met à jour un utilisateur
func (u *userServiceImpl) Update(id int, name string) error {
	// Implementation
	return nil
}

// helper est une fonction privée - pas besoin dans l'interface
func (u *userServiceImpl) helper() {
	// Private method
}

// Config est une struct simple sans méthode - PAS BESOIN D'INTERFACE
type Config struct {
	Host string
	Port int
}

// DataModel est une struct DTO sans méthode - PAS BESOIN D'INTERFACE
type DataModel struct {
	ID   int
	Name string
	Tags []string
}
