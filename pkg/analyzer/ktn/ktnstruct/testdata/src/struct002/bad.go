package struct002

// BadUserService struct sans interface correspondante - VIOLATION
type BadUserService struct { // want "KTN-STRUCT-002"
	users map[int]string
}

// Create crée un utilisateur
func (b *BadUserService) Create(name string) error {
	// Implementation
	return nil
}

// GetByID récupère un utilisateur par ID
func (b *BadUserService) GetByID(id int) (string, error) {
	// Implementation
	return "", nil
}

// IncompleteService interface incomplète (manque Delete)
type IncompleteService interface {
	Save(data string) error
	Load(id int) (string, error)
}

// incompleteServiceImpl a une méthode Delete non dans l'interface - VIOLATION
type incompleteServiceImpl struct { // want "KTN-STRUCT-002"
	data map[int]string
}

// Save est dans l'interface
func (i *incompleteServiceImpl) Save(data string) error {
	// Implementation
	return nil
}

// Load est dans l'interface
func (i *incompleteServiceImpl) Load(id int) (string, error) {
	// Implementation
	return "", nil
}

// Delete est une méthode publique MAIS pas dans l'interface - VIOLATION
func (i *incompleteServiceImpl) Delete(id int) error {
	// Implementation
	return nil
}
