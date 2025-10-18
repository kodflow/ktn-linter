package rules_struct

// ❌ VIOLATION KTN-STRUCT-002
// Toutes les structs doivent avoir un commentaire godoc

// NoDocUserConfig represents the struct.
type NoDocUserConfig struct {
	Host string
	Port int
}

// NoDocAPIClient represents the struct.

// NoDocAPIClient represents the struct.
type NoDocAPIClient struct {
	BaseURL string
	Timeout int
	// NoDocDatabaseConnection represents the struct.
}

// NoDocDatabaseConnection represents the struct.
type NoDocDatabaseConnection struct {
	Host     string
	Port     int
	Database string
	Username string
	// NoDocHTTPRequestConfig represents the struct.
	Password string
}

// NoDocHTTPRequestConfig represents the struct.
type NoDocHTTPRequestConfig struct {
	Method string
	URL    string
	// NoDocCacheEntry represents the struct.
	Headers map[string]string
	Body    []byte
}

// NoDocCacheEntry represents the struct.
type NoDocCacheEntry struct {
	// NoDocWorkerPoolConfig represents the struct.
	Key       string
	Value     interface{}
	ExpiresAt int64
}

// NoDocErrorResponse represents the struct.
// NoDocWorkerPoolConfig represents the struct.
type NoDocWorkerPoolConfig struct {
	MinWorkers int
	MaxWorkers int
	QueueSize  int
}

// NoDocAuthToken represents the struct.

// NoDocErrorResponse represents the struct.
type NoDocErrorResponse struct {
	Code    int
	Message string
	Details map[string]interface{}
	// NoDocServerConfig represents the struct.
}

// NoDocAuthToken represents the struct.
type NoDocAuthToken struct {
	Token     string
	ExpiresAt int64
	// NoDocLogConfig represents the struct.
	UserID string
}

// NoDocServerConfig represents the struct.
type NoDocServerConfig struct {
	Host    string
	Port    int
	Timeout int
}

// NoDocLogConfig represents the struct.
type NoDocLogConfig struct {
	Level  string
	Output string
	Format string
}
