package struct008

// UserEntity entité avec encapsulation correcte - OK
type UserEntity struct {
	id    int
	name  string
	email string
}

// NewUserEntity crée une nouvelle entité utilisateur.
//
// Params:
//   - id: identifiant
//   - name: nom
//   - email: email
//
// Returns:
//   - *UserEntity: instance
func NewUserEntity(id int, name, email string) *UserEntity {
	return &UserEntity{
		id:    id,
		name:  name,
		email: email,
	}
}

// GetID retourne l'identifiant.
//
// Returns:
//   - int: identifiant
func (u *UserEntity) GetID() int {
	return u.id
}

// GetName retourne le nom.
//
// Returns:
//   - string: nom
func (u *UserEntity) GetName() string {
	return u.name
}

// GetEmail retourne l'email.
//
// Returns:
//   - string: email
func (u *UserEntity) GetEmail() string {
	return u.email
}

// Save sauvegarde l'entité
func (u *UserEntity) Save() error {
	return nil
}

// SimpleDTO struct sans méthodes (DTO) - PAS BESOIN DE GETTERS
type SimpleDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Config configuration simple (≤3 champs) - PAS BESOIN DE GETTERS
type Config struct {
	Host string
	Port int
	TLS  bool
}

// internalService service privé - PAS DE RÈGLE
type internalService struct {
	Data string
}

func (i *internalService) process() error {
	return nil
}

// Repository entité avec encapsulation complète.
type Repository struct {
	db     Database
	logger Logger
	cache  Cache
}

// NewRepository crée un nouveau repository.
//
// Params:
//   - db: base de données
//   - logger: logger
//   - cache: cache
//
// Returns:
//   - *Repository: instance
func NewRepository(db Database, logger Logger, cache Cache) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
		cache:  cache,
	}
}

// GetDB retourne la base de données.
//
// Returns:
//   - Database: base
func (r *Repository) GetDB() Database {
	return r.db
}

// GetLogger retourne le logger.
//
// Returns:
//   - Logger: logger
func (r *Repository) GetLogger() Logger {
	return r.logger
}

// GetCache retourne le cache.
//
// Returns:
//   - Cache: cache
func (r *Repository) GetCache() Cache {
	return r.cache
}

// Find recherche une entité
func (r *Repository) Find(id int) (interface{}, error) {
	return nil, nil
}

// EmailService service avec encapsulation partielle OK (≤3 champs).
type EmailService struct {
	host string
	port int
	tls  bool
}

// NewEmailService crée un service email.
//
// Params:
//   - host: hôte
//   - port: port
//   - tls: TLS activé
//
// Returns:
//   - *EmailService: instance
func NewEmailService(host string, port int, tls bool) *EmailService {
	return &EmailService{
		host: host,
		port: port,
		tls:  tls,
	}
}

// Send envoie un email
func (e *EmailService) Send(to, subject, body string) error {
	return nil
}

// Types pour compilation
type Database interface{}
type Logger interface{}
type Cache interface{}
