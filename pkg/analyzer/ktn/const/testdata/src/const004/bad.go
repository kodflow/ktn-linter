package const004

const (
	SMTP_PORT          int     = 25
	SSH_PORT           int     = 22
	TELNET_PORT        int     = 23
	POOL_MAX_SIZE      int     = 500
	POOL_MIN_SIZE      int     = 5
	CONNECTION_TIMEOUT int     = 60
	SERVER_VERSION     string  = "v2.0"
	BASE_PATH          string  = "/base"
	AUTH_TOKEN         string  = "token123"
	AUTO_RELOAD        bool    = false
	STRICT_MODE        bool    = true
	LOG_ENABLED        bool    = true
	ATTEMPTS_LIMIT     int     = 3
	WAIT_TIME_MS       int     = 500
	SCALE_FACTOR       float64 = 2.0
	CACHE_HOST         string  = "127.0.0.1"
	CACHE_PORT         int     = 6379
	SCHEMA_NAME        string  = "public"
	ADMIN_USER         string  = "root"
)
