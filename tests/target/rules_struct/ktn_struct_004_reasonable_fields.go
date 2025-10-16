package rules_struct

// ✅ CONFORME KTN-STRUCT-004
// Les structs ont un nombre raisonnable de champs (≤15)
// Les structs complexes sont décomposées en sous-structs

// SmallStruct a peu de champs (3 ≤ 15).
type SmallStruct struct {
	// Field1 description
	Field1 string
	// Field2 description
	Field2 int
	// Field3 description
	Field3 bool
}

// MediumStruct a un nombre acceptable de champs (10 ≤ 15).
type MediumStruct struct {
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
	Field6 int
	// Field7 description
	Field7 int
	// Field8 description
	Field8 bool
	// Field9 description
	Field9 bool
	// Field10 description
	Field10 []string
}

// MaxStruct est à la limite acceptable (15 ≤ 15).
type MaxStruct struct {
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
	Field11 int
	// Field12 description
	Field12 int
	// Field13 description
	Field13 bool
	// Field14 description
	Field14 bool
	// Field15 description
	Field15 []string
}

// ReasonableServerConfigCore contient la configuration de base du serveur.
type ReasonableServerConfigCore struct {
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
}

// ReasonableServerConfigTLS contient la configuration TLS.
type ReasonableServerConfigTLS struct {
	// Enabled indique si TLS est activé
	Enabled bool
	// CertPath est le chemin du certificat
	CertPath string
	// KeyPath est le chemin de la clé
	KeyPath string
}

// ReasonableServerConfigLogs contient la configuration des logs.
type ReasonableServerConfigLogs struct {
	// Level est le niveau de log
	Level string
	// Output est la sortie des logs
	Output string
}

// ReasonableServerConfigMetrics contient la configuration des métriques.
type ReasonableServerConfigMetrics struct {
	// Enabled indique si les métriques sont activées
	Enabled bool
	// Port est le port des métriques
	Port int
}

// ReasonableServerConfigCORS contient la configuration CORS.
type ReasonableServerConfigCORS struct {
	// Enabled indique si CORS est activé
	Enabled bool
	// AllowedOrigins liste les origines autorisées
	AllowedOrigins []string
	// AllowedMethods liste les méthodes autorisées
	AllowedMethods []string
}

// ReasonableServerConfig contient toute la configuration du serveur décomposée.
type ReasonableServerConfig struct {
	// Core contient la configuration de base
	Core ReasonableServerConfigCore
	// TLS contient la configuration TLS
	TLS ReasonableServerConfigTLS
	// Logs contient la configuration des logs
	Logs ReasonableServerConfigLogs
	// Metrics contient la configuration des métriques
	Metrics ReasonableServerConfigMetrics
	// CORS contient la configuration CORS
	CORS ReasonableServerConfigCORS
	// HealthCheckPath est le chemin du health check
	HealthCheckPath string
	// ShutdownTimeout est le délai d'arrêt
	ShutdownTimeout int
	// MaxBodySize est la taille maximale du body
	MaxBodySize int64
}

// UserBasicInfo contient les informations de base d'un utilisateur.
type UserBasicInfo struct {
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
}

// UserContactInfo contient les informations de contact.
type UserContactInfo struct {
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
}

// UserMetadata contient les métadonnées de l'utilisateur.
type UserMetadata struct {
	// BirthDate est la date de naissance
	BirthDate string
	// CreatedAt est la date de création
	CreatedAt int64
	// UpdatedAt est la date de mise à jour
	UpdatedAt int64
	// LastLoginAt est la date de dernière connexion
	LastLoginAt int64
}

// UserPermissions contient les permissions de l'utilisateur.
type UserPermissions struct {
	// IsActive indique si l'utilisateur est actif
	IsActive bool
	// IsVerified indique si l'email est vérifié
	IsVerified bool
	// Roles liste les rôles de l'utilisateur
	Roles []string
	// Permissions liste les permissions
	Permissions []string
}

// User représente un utilisateur complet avec toutes ses informations décomposées.
type User struct {
	// BasicInfo contient les informations de base
	BasicInfo UserBasicInfo
	// ContactInfo contient les informations de contact
	ContactInfo UserContactInfo
	// Metadata contient les métadonnées
	Metadata UserMetadata
	// Permissions contient les permissions
	Permissions UserPermissions
}
