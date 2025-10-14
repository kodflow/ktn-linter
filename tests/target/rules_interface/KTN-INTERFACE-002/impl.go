package KTN_INTERFACE_002

// userServiceImplI002Good est l'implémentation privée.
type userServiceImplI002Good struct {
	db       string
	cacheDir string
}

// ProcessI002Good traite des données.
//
// Params:
//   - data: les données à traiter
//
// Returns:
//   - error: une erreur si l'opération échoue
func (s *userServiceImplI002Good) ProcessI002Good(data string) error {
	return nil
}

// NewUserServiceI002Good crée une instance.
//
// Params:
//   - db: base de données
//   - cacheDir: répertoire cache
//
// Returns:
//   - UserServiceI002Good: nouvelle instance
func NewUserServiceI002Good(db string, cacheDir string) UserServiceI002Good {
	return &userServiceImplI002Good{db: db, cacheDir: cacheDir}
}

// orderManagerImplI002Good est l'implémentation privée.
type orderManagerImplI002Good struct {
	orders []string
}

// CreateOrderI002Good crée une commande.
//
// Params:
//   - order: la commande à créer
//
// Returns:
//   - error: une erreur si l'opération échoue
func (o *orderManagerImplI002Good) CreateOrderI002Good(order string) error {
	o.orders = append(o.orders, order)
	return nil
}

// NewOrderManagerI002Good crée une instance.
//
// Returns:
//   - OrderManagerI002Good: nouvelle instance
func NewOrderManagerI002Good() OrderManagerI002Good {
	return &orderManagerImplI002Good{orders: []string{}}
}
