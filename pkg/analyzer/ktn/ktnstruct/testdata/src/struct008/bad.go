package struct008

// BadUser champs publics alors qu'il a des méthodes - VIOLATION
type BadUser struct { // want "KTN-STRUCT-008"
	ID    int
	Name  string
	Email string
	Role  string
}

// Save méthode présente mais champs publics
func (b *BadUser) Save() error {
	return nil
}

// Delete méthode présente
func (b *BadUser) Delete() error {
	return nil
}

// MixedVisibility mélange champs publics/privés avec méthodes - VIOLATION
type MixedVisibility struct { // want "KTN-STRUCT-008"
	ID       int
	name     string
	Email    string
	password string
}

// Validate méthode présente
func (m *MixedVisibility) Validate() bool {
	return true
}

// ProductEntity entité avec champs publics - VIOLATION
type ProductEntity struct { // want "KTN-STRUCT-008"
	ID          int
	Name        string
	Description string
	Price       float64
}

// GetFormattedPrice méthode présente
func (p *ProductEntity) GetFormattedPrice() string {
	return ""
}

// UpdatePrice méthode présente
func (p *ProductEntity) UpdatePrice(newPrice float64) {
	p.Price = newPrice
}

// NoGetters champs privés MAIS pas de getters - VIOLATION
type NoGetters struct { // want "KTN-STRUCT-008"
	id    int
	name  string
	email string
	role  string
}

// Save méthode présente
func (n *NoGetters) Save() error {
	return nil
}

// Update méthode présente mais pas de getters pour les champs
func (n *NoGetters) Update(name string) {
	n.name = name
}

// PartialGetters champs privés avec getters incomplets - VIOLATION
type PartialGetters struct { // want "KTN-STRUCT-008"
	id    int
	name  string
	email string
	role  string
}

// GetID getter présent
func (p *PartialGetters) GetID() int {
	return p.id
}

// GetName getter présent
func (p *PartialGetters) GetName() string {
	return p.name
}

// Manque GetEmail et GetRole - VIOLATION

// Process méthode présente
func (p *PartialGetters) Process() error {
	return nil
}
