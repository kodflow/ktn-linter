package rules_struct

// ✅ CONFORME KTN-STRUCT-004
// Les structs ont un nombre raisonnable de champs (≤15)
// Les structs complexes sont décomposées en sous-structs

// smallStruct a peu de champs (3 ≤ 15).
type smallStruct struct {
	// Field1 description
	Field1 string
	// Field2 description
	Field2 int
	// Field3 description
	Field3 bool
}

// mediumStruct a un nombre acceptable de champs (10 ≤ 15).
type mediumStruct struct {
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

// maxStruct est à la limite acceptable (15 ≤ 15).
type maxStruct struct {
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

// reasonableServerConfigCore contient la configuration de base du serveur.
type reasonableServerConfigCore struct {
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

// reasonableServerConfigTLS contient la configuration TLS.
type reasonableServerConfigTLS struct {
	// Enabled indique si TLS est activé
	Enabled bool
	// CertPath est le chemin du certificat
	CertPath string
	// KeyPath est le chemin de la clé
	KeyPath string
}

// reasonableServerConfigLogs contient la configuration des logs.
type reasonableServerConfigLogs struct {
	// Level est le niveau de log
	Level string
	// Output est la sortie des logs
	Output string
}

// reasonableServerConfigMetrics contient la configuration des métriques.
type reasonableServerConfigMetrics struct {
	// Enabled indique si les métriques sont activées
	Enabled bool
	// Port est le port des métriques
	Port int
}

// reasonableServerConfigCORS contient la configuration CORS.
type reasonableServerConfigCORS struct {
	// Enabled indique si CORS est activé
	Enabled bool
	// AllowedOrigins liste les origines autorisées
	AllowedOrigins []string
	// AllowedMethods liste les méthodes autorisées
	AllowedMethods []string
}

// reasonableServerConfig contient toute la configuration du serveur décomposée.
type reasonableServerConfig struct {
	// Core contient la configuration de base
	Core reasonableServerConfigCore
	// TLS contient la configuration TLS
	TLS reasonableServerConfigTLS
	// Logs contient la configuration des logs
	Logs reasonableServerConfigLogs
	// Metrics contient la configuration des métriques
	Metrics reasonableServerConfigMetrics
	// CORS contient la configuration CORS
	CORS reasonableServerConfigCORS
	// HealthCheckPath est le chemin du health check
	HealthCheckPath string
	// ShutdownTimeout est le délai d'arrêt
	ShutdownTimeout int
	// MaxBodySize est la taille maximale du body
	MaxBodySize int64
}

// userBasicInfo contient les informations de base d'un utilisateur.
type userBasicInfo struct {
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

// userContactInfo contient les informations de contact.
type userContactInfo struct {
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

// userMetadata contient les métadonnées de l'utilisateur.
type userMetadata struct {
	// BirthDate est la date de naissance
	BirthDate string
	// CreatedAt est la date de création
	CreatedAt int64
	// UpdatedAt est la date de mise à jour
	UpdatedAt int64
	// LastLoginAt est la date de dernière connexion
	LastLoginAt int64
}

// userPermissions contient les permissions de l'utilisateur.
type userPermissions struct {
	// IsActive indique si l'utilisateur est actif
	IsActive bool
	// IsVerified indique si l'email est vérifié
	IsVerified bool
	// Roles liste les rôles de l'utilisateur
	Roles []string
	// Permissions liste les permissions
	Permissions []string
}

// user représente un utilisateur complet avec toutes ses informations décomposées.
type user struct {
	// BasicInfo contient les informations de base
	BasicInfo userBasicInfo
	// ContactInfo contient les informations de contact
	ContactInfo userContactInfo
	// Metadata contient les métadonnées
	Metadata userMetadata
	// Permissions contient les permissions
	Permissions userPermissions
}
