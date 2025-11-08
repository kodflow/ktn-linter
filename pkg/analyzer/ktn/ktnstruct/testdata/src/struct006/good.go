package struct006

// UserEntity entité avec encapsulation correcte - OK.
// Représente une entité utilisateur avec champs privés et getters publics.
type UserEntity struct {
	id    int
	name  string
	email string
}

// UserEntityInterface définit les méthodes publiques de UserEntity.
type UserEntityInterface interface {
	GetID() int
	GetName() string
	GetEmail() string
	Save() error
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
	// Retourne nouvelle instance avec champs initialisés
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
	// Retourne le champ id
	return u.id
}

// GetName retourne le nom.
//
// Returns:
//   - string: nom
func (u *UserEntity) GetName() string {
	// Retourne le champ name
	return u.name
}

// GetEmail retourne l'email.
//
// Returns:
//   - string: email
func (u *UserEntity) GetEmail() string {
	// Retourne le champ email
	return u.email
}

// Save sauvegarde l'entité.
//
// Returns:
//   - error: erreur éventuelle
func (u *UserEntity) Save() error {
	// Retourne nil si succès
	return nil
}

// SimpleDTO struct sans méthodes (DTO) - PAS BESOIN DE GETTERS.
// Structure de transfert de données simple.
type SimpleDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Config configuration simple (≤3 champs) - PAS BESOIN DE GETTERS.
// Contient les paramètres de configuration réseau.
type Config struct {
	Host string
	Port int
	TLS  bool
}

// internalService service privé - PAS DE RÈGLE
type internalService struct {
	Data string
}

// process traite les données.
//
// Returns:
//   - error: erreur éventuelle
func (i *internalService) process() error {
	// Retourne nil si succès
	return nil
}

// Repository entité avec encapsulation complète.
// Service de persistance avec dépendances et méthodes publiques.
type Repository struct {
	db     Database
	logger Logger
	cache  Cache
}

// RepositoryInterface définit les méthodes publiques de Repository.
type RepositoryInterface interface {
	GetDB() Database
	GetLogger() Logger
	GetCache() Cache
	Find(id int) (interface{}, error)
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
	// Retourne nouvelle instance avec dépendances
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
	// Retourne le champ db
	return r.db
}

// GetLogger retourne le logger.
//
// Returns:
//   - Logger: logger
func (r *Repository) GetLogger() Logger {
	// Retourne le champ logger
	return r.logger
}

// GetCache retourne le cache.
//
// Returns:
//   - Cache: cache
func (r *Repository) GetCache() Cache {
	// Retourne le champ cache
	return r.cache
}

// Find recherche une entité.
//
// Params:
//   - id: identifiant de l'entité
//
// Returns:
//   - interface{}: entité trouvée
//   - error: erreur éventuelle
func (r *Repository) Find(id int) (interface{}, error) {
	// Retourne nil si non trouvé
	return nil, nil
}

// EmailService service avec encapsulation partielle OK (≤3 champs).
// Service d'envoi d'emails avec configuration réseau.
type EmailService struct {
	host string
	port int
	tls  bool
}

// EmailServiceInterface définit les méthodes de EmailService.
type EmailServiceInterface interface {
	Send(to, subject, body string) error
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
	// Retourne nouvelle instance avec configuration
	return &EmailService{
		host: host,
		port: port,
		tls:  tls,
	}
}

// Send envoie un email.
//
// Params:
//   - to: destinataire
//   - subject: sujet
//   - body: corps du message
//
// Returns:
//   - error: erreur éventuelle
func (e *EmailService) Send(to, subject, body string) error {
	// Retourne nil si succès
	return nil
}

// Types pour compilation
type Database interface{}
type Logger interface{}
type Cache interface{}
