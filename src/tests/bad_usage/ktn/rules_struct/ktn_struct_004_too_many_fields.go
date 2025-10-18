package rules_struct

// ❌ VIOLATION KTN-STRUCT-004
// Les structs ne doivent pas avoir plus de 15 champs

// LargeStruct a trop de champs (16 > 15).
type LargeStruct struct {
	// Field1 description
	Field1 string
	// Field2 description
	Field2 string
	// Field3 description
	Field3 string
	// Field4 description
	Field4 string
	// Field5 description
	Field5 string
	// Field6 description
	Field6 string
	// Field7 description
	Field7 string
	// Field8 description
	Field8 string
	// Field9 description
	Field9 string
	// Field10 description
	Field10 string
	// Field11 description
	Field11 string
	// Field12 description
	Field12 string
	// Field13 description
	Field13 string
	// Field14 description
	Field14 string
	// Field15 description
	Field15 string
	// Field16 description
	Field16 string
}

// MassiveConfig a beaucoup trop de champs (20 > 15).
type MassiveConfig struct {
	// Host est l'hôte du serveur
	Host string
	// Port est le port d'écoute
	Port int
	// Timeout est le délai d'expiration
	Timeout int
	// MaxConnections est le nombre maximum de connexions
	MaxConnections int
	// ReadTimeout est le délai de lecture
	ReadTimeout int
	// WriteTimeout est le délai d'écriture
	WriteTimeout int
	// IdleTimeout est le délai d'inactivité
	IdleTimeout int
	// TLSEnabled indique si TLS est activé
	TLSEnabled bool
	// CertPath est le chemin du certificat
	CertPath string
	// KeyPath est le chemin de la clé
	KeyPath string
	// LogLevel est le niveau de log
	LogLevel string
	// LogOutput est la sortie des logs
	LogOutput string
	// MetricsEnabled indique si les métriques sont activées
	MetricsEnabled bool
	// MetricsPort est le port des métriques
	MetricsPort int
	// HealthCheckPath est le chemin du health check
	HealthCheckPath string
	// ShutdownTimeout est le délai d'arrêt
	ShutdownTimeout int
	// MaxBodySize est la taille maximale du body
	MaxBodySize int64
	// EnableCORS indique si CORS est activé
	EnableCORS bool
	// AllowedOrigins liste les origines autorisées
	AllowedOrigins []string
	// AllowedMethods liste les méthodes autorisées
	AllowedMethods []string
}

// ComplexUser représente un utilisateur avec trop d'informations (18 > 15).
type ComplexUser struct {
	// ID est l'identifiant unique
	ID string
	// Username est le nom d'utilisateur
	Username string
	// Email est l'adresse email
	Email string
	// FirstName est le prénom
	FirstName string
	// LastName est le nom de famille
	LastName string
	// PhoneNumber est le numéro de téléphone
	PhoneNumber string
	// Address est l'adresse
	Address string
	// City est la ville
	City string
	// ZipCode est le code postal
	ZipCode string
	// Country est le pays
	Country string
	// BirthDate est la date de naissance
	BirthDate string
	// CreatedAt est la date de création
	CreatedAt int64
	// UpdatedAt est la date de mise à jour
	UpdatedAt int64
	// LastLoginAt est la date de dernière connexion
	LastLoginAt int64
	// IsActive indique si l'utilisateur est actif
	IsActive bool
	// IsVerified indique si l'email est vérifié
	IsVerified bool
	// Roles liste les rôles de l'utilisateur
	Roles []string
	// Permissions liste les permissions
	Permissions []string
}
