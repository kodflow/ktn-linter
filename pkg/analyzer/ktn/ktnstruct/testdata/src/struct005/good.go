package struct005

// UserService gère les utilisateurs du système.
// Encapsule la logique métier liée aux utilisateurs.
type UserService struct {
	users map[int]string
}

// NewUserService crée un nouveau service utilisateur.
//
// Returns:
//   - *UserService: instance du service
func NewUserService() *UserService {
	return &UserService{
		users: make(map[int]string),
	}
}

// Create crée un utilisateur
func (u *UserService) Create(name string) error {
	return nil
}

// GetByID récupère un utilisateur
func (u *UserService) GetByID(id int) string {
	return ""
}

// Config représente la configuration.
// Simple DTO sans comportement - PAS DE CONSTRUCTEUR REQUIS.
type Config struct {
	Host string
	Port int
}

// UserDTO représente un utilisateur pour le transfert.
// Pas de logique métier, juste des données - PAS DE CONSTRUCTEUR REQUIS.
type UserDTO struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// internalCache cache privé avec méthodes - PAS DE CONSTRUCTEUR REQUIS.
type internalCache struct {
	data map[string]interface{}
}

func (c *internalCache) get(key string) interface{} {
	return c.data[key]
}

// Repository gère la persistance des données.
// Service avec dépendances.
type Repository struct {
	db     Database
	logger Logger
}

// NewRepository crée un nouveau repository.
//
// Params:
//   - db: instance de base de données
//   - logger: logger pour les traces
//
// Returns:
//   - *Repository: instance du repository
func NewRepository(db Database, logger Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Save sauvegarde une entité
func (r *Repository) Save(entity interface{}) error {
	return nil
}

// EmailService gère l'envoi d'emails.
// Service avec configuration complexe.
type EmailService struct {
	host     string
	port     int
	username string
	password string
	tls      bool
}

// EmailServiceConfig configuration pour EmailService.
type EmailServiceConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	TLS      bool
}

// NewEmailService crée un service email avec config.
//
// Params:
//   - cfg: configuration du service
//
// Returns:
//   - *EmailService: instance du service
func NewEmailService(cfg EmailServiceConfig) *EmailService {
	return &EmailService{
		host:     cfg.Host,
		port:     cfg.Port,
		username: cfg.Username,
		password: cfg.Password,
		tls:      cfg.TLS,
	}
}

// Send envoie un email
func (e *EmailService) Send(to, subject, body string) error {
	return nil
}

// Validator valide des données.
// Service retournant valeur directe (pas pointeur).
type Validator struct {
	rules map[string]func(string) bool
}

// NewValidator crée un validateur - RETOURNE VALEUR
//
// Returns:
//   - Validator: instance du validateur
func NewValidator() Validator {
	return Validator{
		rules: make(map[string]func(string) bool),
	}
}

// Validate valide une donnée
func (v Validator) Validate(key, value string) bool {
	return true
}

// Types pour compilation
type Database interface{}
type Logger interface{}
