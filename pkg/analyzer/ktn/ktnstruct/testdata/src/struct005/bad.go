package struct005

// BadUserService service sans constructeur - VIOLATION
type BadUserService struct { // want "KTN-STRUCT-005"
	users map[int]string
}

// Create méthode présente mais pas de NewBadUserService()
func (b *BadUserService) Create(name string) error {
	return nil
}

// GetByID méthode présente
func (b *BadUserService) GetByID(id int) string {
	return ""
}

// MisnamedService service avec constructeur mal nommé - VIOLATION
type MisnamedService struct { // want "KTN-STRUCT-005"
	data map[string]string
}

// CreateService mauvais nom (devrait être NewMisnamedService) - NE COMPTE PAS
func CreateService() *MisnamedService {
	return &MisnamedService{data: make(map[string]string)}
}

// Process méthode présente
func (m *MisnamedService) Process(key string) error {
	return nil
}

// Cache gère un cache - VIOLATION
type Cache struct { // want "KTN-STRUCT-005"
	items map[string]interface{}
}

// Get récupère un item
func (c *Cache) Get(key string) interface{} {
	return nil
}

// Set définit un item
func (c *Cache) Set(key string, value interface{}) {
	c.items[key] = value
}

// Delete supprime un item
func (c *Cache) Delete(key string) {
	delete(c.items, key)
}

// Clear vide le cache
func (c *Cache) Clear() {
	c.items = make(map[string]interface{})
}

// WrongReturnType constructeur retournant mauvais type - VIOLATION
type WrongReturnType struct { // want "KTN-STRUCT-005"
	value int
}

// NewWrongReturnType retourne string au lieu de *WrongReturnType - NE COMPTE PAS
func NewWrongReturnType() string {
	return "wrong"
}

// GetValue méthode présente
func (w *WrongReturnType) GetValue() int {
	return w.value
}
