package rules_struct

// ❌ VIOLATION KTN-STRUCT-003
// Tous les champs exportés doivent être documentés

// NoCommentUserConfig contient la configuration utilisateur.
type NoCommentUserConfig struct {
	Host string
	Port int
}

// NoCommentAPIClient représente un client API.
type NoCommentAPIClient struct {
	BaseURL    string
	Timeout    int
	MaxRetries int
}

// NoCommentDatabaseConnection représente une connexion à la base de données.
type NoCommentDatabaseConnection struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

// NoCommentHTTPRequestConfig contient la configuration d'une requête HTTP.
type NoCommentHTTPRequestConfig struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
	Timeout int
}

// NoCommentCacheEntry représente une entrée dans le cache.
type NoCommentCacheEntry struct {
	Key       string
	Value     interface{}
	ExpiresAt int64
	Priority  int
}

// NoCommentWorkerPoolConfig contient la configuration d'un pool de workers.
type NoCommentWorkerPoolConfig struct {
	MinWorkers  int
	MaxWorkers  int
	QueueSize   int
	IdleTimeout int
}

// NoCommentErrorResponse représente une réponse d'erreur.
type NoCommentErrorResponse struct {
	Code      int
	Message   string
	Details   map[string]interface{}
	Timestamp int64
}

// NoCommentAuthToken représente un token d'authentification.
type NoCommentAuthToken struct {
	Token     string
	ExpiresAt int64
	UserID    string
	Scope     []string
}

// NoCommentServerConfig contient la configuration du serveur.
type NoCommentServerConfig struct {
	Host     string
	Port     int
	Timeout  int
	TLS      bool
	CertPath string
}

// NoCommentLogConfig contient la configuration des logs.
type NoCommentLogConfig struct {
	Level   string
	Output  string
	Format  string
	MaxSize int
	MaxAge  int
}
