package rules_struct

// ❌ VIOLATION KTN-STRUCT-002
// Toutes les structs doivent avoir un commentaire godoc

type NoDocUserConfig struct {
	Host string
	Port int
}

type NoDocAPIClient struct {
	BaseURL string
	Timeout int
}

type NoDocDatabaseConnection struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

type NoDocHTTPRequestConfig struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
}

type NoDocCacheEntry struct {
	Key       string
	Value     interface{}
	ExpiresAt int64
}

type NoDocWorkerPoolConfig struct {
	MinWorkers int
	MaxWorkers int
	QueueSize  int
}

type NoDocErrorResponse struct {
	Code    int
	Message string
	Details map[string]interface{}
}

type NoDocAuthToken struct {
	Token     string
	ExpiresAt int64
	UserID    string
}

type NoDocServerConfig struct {
	Host    string
	Port    int
	Timeout int
}

type NoDocLogConfig struct {
	Level  string
	Output string
	Format string
}
