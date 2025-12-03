package struct006

// GoodUser struct avec getters idiomatiques Go (sans préfixe Get)
type GoodUser struct {
	id    int
	name  string
	email string
}

// ID retourne l'identifiant - OK (convention Go idiomatique)
func (u *GoodUser) ID() int {
	// Retourne le champ id
	return u.id
}

// Name retourne le nom - OK (convention Go idiomatique)
func (u *GoodUser) Name() string {
	// Retourne le champ name
	return u.name
}

// Email retourne l'email - OK (convention Go idiomatique)
func (u *GoodUser) Email() string {
	// Retourne le champ email
	return u.email
}

// SetName définit le nom - OK (setter garde le préfixe Set)
func (u *GoodUser) SetName(name string) {
	// Modifie le champ name
	u.name = name
}

// GoodConfig struct de configuration avec getters idiomatiques
type GoodConfig struct {
	host string
	port int
	tls  bool
}

// Host retourne l'hôte - OK (convention Go idiomatique)
func (c *GoodConfig) Host() string {
	// Retourne le champ host
	return c.host
}

// Port retourne le port - OK (convention Go idiomatique)
func (c *GoodConfig) Port() int {
	// Retourne le champ port
	return c.port
}

// TLS retourne si TLS est activé - OK (convention Go idiomatique)
func (c *GoodConfig) TLS() bool {
	// Retourne le champ tls
	return c.tls
}

// GoodRepository struct avec getters idiomatiques
type GoodRepository struct {
	db     GoodDatabase
	logger GoodLogger
}

// DB retourne la base de données - OK (convention Go idiomatique)
func (r *GoodRepository) DB() GoodDatabase {
	// Retourne le champ db
	return r.db
}

// Logger retourne le logger - OK (convention Go idiomatique)
func (r *GoodRepository) Logger() GoodLogger {
	// Retourne le champ logger
	return r.logger
}

// Find recherche une entité - OK (pas un getter)
func (r *GoodRepository) Find(id int) (interface{}, error) {
	// Retourne nil si non trouvé
	return nil, nil
}

// GetByID est une méthode de recherche, pas un getter - OK
// Elle prend un paramètre, donc ce n'est pas un getter simple
func (r *GoodRepository) GetByID(id int) (interface{}, error) {
	// Retourne nil si non trouvé
	return nil, nil
}

// privateStruct struct privée - PAS DE VÉRIFICATION
type privateStruct struct {
	value string
}

// GetValue méthode sur struct privée - OK (struct non exportée)
func (p *privateStruct) GetValue() string {
	// Retourne le champ value
	return p.value
}

// Types pour compilation
type GoodDatabase interface{}
type GoodLogger interface{}
