package rules_struct

// ❌ VIOLATION KTN-STRUCT-001
// Les structs doivent utiliser MixedCaps, pas snake_case

// user_config contient la configuration utilisateur.
type user_config struct {
	Host string
	Port int
}

// api_client représente un client API.
type api_client struct {
	BaseURL string
	Timeout int
}

// database_connection représente une connexion à la base de données.
type database_connection struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

// http_request_config contient la configuration d'une requête HTTP.
type http_request_config struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
}

// cache_entry représente une entrée dans le cache.
type cache_entry struct {
	Key       string
	Value     interface{}
	ExpiresAt int64
}

// worker_pool_config contient la configuration d'un pool de workers.
type worker_pool_config struct {
	MinWorkers int
	MaxWorkers int
	QueueSize  int
}

// error_response représente une réponse d'erreur.
type error_response struct {
	Code    int
	Message string
	Details map[string]interface{}
}

// auth_token représente un token d'authentification.
type auth_token struct {
	Token     string
	ExpiresAt int64
	UserID    string
}
