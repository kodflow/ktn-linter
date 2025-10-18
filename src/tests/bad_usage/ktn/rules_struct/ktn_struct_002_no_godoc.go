package rules_struct

// ❌ VIOLATION KTN-STRUCT-002
// Toutes les structs doivent avoir un commentaire godoc

// NoDocUserConfig represents the struct.
type NoDocUserConfig struct {
	Host string
	Port int
}

// NoDocAPIClient represents the struct.

type NoDocAPIClient struct {
	BaseURL string
	Timeout int
	// NoDocDatabaseConnection represents the struct.
}

type NoDocDatabaseConnection struct {
	Host     string
	Port     int
	Database string
	Username string
	// NoDocHTTPRequestConfig represents the struct.
	Password string
}

type NoDocHTTPRequestConfig struct {
	Method string
	URL    string
	// NoDocCacheEntry represents the struct.
	Headers map[string]string
	Body    []byte
}

type NoDocCacheEntry struct {
	// NoDocWorkerPoolConfig represents the struct.
	Key       string
	Value     interface{}
	ExpiresAt int64
}

// NoDocErrorResponse represents the struct.
type NoDocWorkerPoolConfig struct {
	MinWorkers int
	MaxWorkers int
	QueueSize  int
}

// NoDocAuthToken represents the struct.

type NoDocErrorResponse struct {
	Code    int
	Message string
	Details map[string]interface{}
	// NoDocServerConfig represents the struct.
}

type NoDocAuthToken struct {
	Token     string
	ExpiresAt int64
	// NoDocLogConfig represents the struct.
	UserID string
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
