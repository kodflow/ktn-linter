package struct006

// User struct avec getters utilisant le préfixe Get (non idiomatique)
type User struct {
	id    int
	name  string
	email string
}

// GetID retourne l'identifiant - VIOLATION: devrait être ID()
func (u *User) GetID() int { // want "KTN-STRUCT-006"
	// Retourne le champ id
	return u.id
}

// GetName retourne le nom - VIOLATION: devrait être Name()
func (u *User) GetName() string { // want "KTN-STRUCT-006"
	// Retourne le champ name
	return u.name
}

// GetEmail retourne l'email - VIOLATION: devrait être Email()
func (u *User) GetEmail() string { // want "KTN-STRUCT-006"
	// Retourne le champ email
	return u.email
}

// Save sauvegarde l'utilisateur - OK (pas un getter)
func (u *User) Save() error {
	// Retourne nil si succès
	return nil
}

// Config struct de configuration avec getters
type Config struct {
	host string
	port int
	tls  bool
}

// GetHost retourne l'hôte - VIOLATION: devrait être Host()
func (c *Config) GetHost() string { // want "KTN-STRUCT-006"
	// Retourne le champ host
	return c.host
}

// GetPort retourne le port - VIOLATION: devrait être Port()
func (c *Config) GetPort() int { // want "KTN-STRUCT-006"
	// Retourne le champ port
	return c.port
}

// Repository struct avec getters
type Repository struct {
	db     Database
	logger Logger
}

// GetDB retourne la base de données - VIOLATION: devrait être DB()
func (r *Repository) GetDB() Database { // want "KTN-STRUCT-006"
	// Retourne le champ db
	return r.db
}

// GetLogger retourne le logger - VIOLATION: devrait être Logger()
func (r *Repository) GetLogger() Logger { // want "KTN-STRUCT-006"
	// Retourne le champ logger
	return r.logger
}

// Types pour compilation
type Database interface{}
type Logger interface{}
