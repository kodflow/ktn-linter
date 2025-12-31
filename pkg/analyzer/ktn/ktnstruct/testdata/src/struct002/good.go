// Package struct002 provides good test cases.
package struct002

const (
	// defaultMapSize taille par défaut des maps
	defaultMapSize int = 10
)

// UserServiceConfig gère les utilisateurs du système.
// Encapsule la logique métier liée aux utilisateurs.
type UserServiceConfig struct {
	users map[int]string
}

// NewUserService crée un nouveau service utilisateur.
//
// Returns:
//   - *UserServiceConfig: instance du service
func NewUserService() *UserServiceConfig {
	// Retourne nouvelle instance avec map initialisée
	return &UserServiceConfig{
		users: make(map[int]string, defaultMapSize),
	}
}

// Create crée un utilisateur.
//
// Params:
//   - _name: nom de l'utilisateur (non utilisé dans cette implémentation)
//
// Returns:
//   - error: erreur éventuelle
func (u *UserServiceConfig) Create(_name string) error {
	// Retourne nil si succès
	return nil
}

// GetByID récupère un utilisateur.
//
// Params:
//   - _id: identifiant de l'utilisateur (non utilisé dans cette implémentation)
//
// Returns:
//   - string: nom de l'utilisateur
func (u *UserServiceConfig) GetByID(_id int) string {
	// Retourne chaîne vide si non trouvé
	return ""
}

// Config représente la configuration.
// Simple DTO sans comportement - PAS DE CONSTRUCTEUR REQUIS.
type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// UserDTO représente un utilisateur pour le transfert.
// Pas de logique métier, juste des données - PAS DE CONSTRUCTEUR REQUIS.
type UserDTO struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// internalCacheData cache privé avec méthodes - PAS DE CONSTRUCTEUR REQUIS.
type internalCacheData struct {
	data map[string]any
}

// RepositoryConfig gère la persistance des données.
// Service avec dépendances.
type RepositoryConfig struct {
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
//   - *RepositoryConfig: instance du repository
func NewRepository(db Database, logger Logger) *RepositoryConfig {
	// Retourne nouvelle instance du repository
	return &RepositoryConfig{
		db:     db,
		logger: logger,
	}
}

// Save sauvegarde une entité.
//
// Params:
//   - _entity: entité à sauvegarder (non utilisée dans cette implémentation)
//
// Returns:
//   - error: erreur éventuelle
func (r *RepositoryConfig) Save(_entity any) error {
	// Retourne nil si succès
	return nil
}

// EmailServiceSettings gère l'envoi d'emails.
// Service avec configuration complexe.
type EmailServiceSettings struct {
	host     string
	port     int
	username string
	password string
	tls      bool
}

// EmailServiceConfig configuration pour EmailService.
// Contient tous les paramètres nécessaires à la création du service email.
type EmailServiceConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	TLS      bool   `json:"tls"`
}

// NewEmailService crée un service email avec config.
//
// Params:
//   - cfg: configuration du service
//
// Returns:
//   - *EmailServiceSettings: instance du service
func NewEmailService(cfg EmailServiceConfig) *EmailServiceSettings {
	// Retourne nouvelle instance avec configuration
	return &EmailServiceSettings{
		host:     cfg.Host,
		port:     cfg.Port,
		username: cfg.Username,
		password: cfg.Password,
		tls:      cfg.TLS,
	}
}

// Host retourne l'hôte.
//
// Returns:
//   - string: hôte du serveur
func (e *EmailServiceSettings) Host() string {
	// Retourne le champ host
	return e.host
}

// Port retourne le port.
//
// Returns:
//   - int: port du serveur
func (e *EmailServiceSettings) Port() int {
	// Retourne le champ port
	return e.port
}

// Username retourne le nom d'utilisateur.
//
// Returns:
//   - string: nom d'utilisateur
func (e *EmailServiceSettings) Username() string {
	// Retourne le champ username
	return e.username
}

// Password retourne le mot de passe.
//
// Returns:
//   - string: mot de passe
func (e *EmailServiceSettings) Password() string {
	// Retourne le champ password
	return e.password
}

// Tls retourne l'état TLS.
//
// Returns:
//   - bool: true si TLS activé
func (e *EmailServiceSettings) Tls() bool {
	// Retourne le champ tls
	return e.tls
}

// Send envoie un email.
//
// Params:
//   - _to: destinataire (non utilisé dans cette implémentation)
//   - _subject: sujet du message (non utilisé dans cette implémentation)
//   - _body: corps du message (non utilisé dans cette implémentation)
//
// Returns:
//   - error: erreur éventuelle
func (e *EmailServiceSettings) Send(_to, _subject, _body string) error {
	// Retourne nil si succès
	return nil
}

// ValidatorConfig valide des données.
// Service retournant valeur directe (pas pointeur).
type ValidatorConfig struct {
	rules map[string]func(string) bool
}

// NewValidator crée un validateur - RETOURNE VALEUR
//
// Returns:
//   - ValidatorConfig: instance du validateur
func NewValidator() ValidatorConfig {
	// Retourne nouvelle instance avec map initialisée
	return ValidatorConfig{
		rules: make(map[string]func(string) bool, defaultMapSize),
	}
}

// Validate valide une donnée.
//
// Params:
//   - _key: clé de la règle (non utilisée dans cette implémentation)
//   - _value: valeur à valider (non utilisée dans cette implémentation)
//
// Returns:
//   - bool: true si valide
func (v ValidatorConfig) Validate(_key, _value string) bool {
	// Retourne true si valide
	return true
}

// Database interface pour la base de données.
type Database any

// Logger interface pour les logs.
type Logger any
