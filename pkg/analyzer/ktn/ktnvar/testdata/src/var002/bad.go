// Bad examples for the var002 test case.
package var002

const (
	// SMTP_PORT_VALUE defines SMTP port value
	SMTP_PORT_VALUE int = 25
	// SSH_PORT_VALUE defines SSH port value
	SSH_PORT_VALUE int = 22
	// TELNET_PORT_VALUE defines telnet port value
	TELNET_PORT_VALUE int = 23
	// POOL_MAX_SIZE_VALUE defines pool max size
	POOL_MAX_SIZE_VALUE int = 500
	// POOL_MIN_SIZE_VALUE defines pool min size
	POOL_MIN_SIZE_VALUE int = 5
	// CONNECTION_TIMEOUT_VALUE defines connection timeout
	CONNECTION_TIMEOUT_VALUE int = 60
	// SERVER_VERSION_VALUE defines server version
	SERVER_VERSION_VALUE string = "v2.0"
	// BASE_PATH_VALUE defines base path
	BASE_PATH_VALUE string = "/base"
	// AUTH_TOKEN_VALUE defines auth token
	AUTH_TOKEN_VALUE string = "token123"
	// AUTO_RELOAD_VALUE defines auto reload setting
	AUTO_RELOAD_VALUE bool = false
	// STRICT_MODE_VALUE defines strict mode setting
	STRICT_MODE_VALUE bool = true
	// LOG_ENABLED_VALUE defines log enabled setting
	LOG_ENABLED_VALUE bool = true
	// ATTEMPTS_LIMIT_VALUE defines attempts limit
	ATTEMPTS_LIMIT_VALUE int = 3
	// WAIT_TIME_MS_VALUE defines wait time in ms
	WAIT_TIME_MS_VALUE int = 500
	// SCALE_FACTOR_VALUE defines scale factor
	SCALE_FACTOR_VALUE float64 = 2.0
	// CACHE_HOST_VALUE defines cache host
	CACHE_HOST_VALUE string = "127.0.0.1"
	// CACHE_PORT_VALUE defines cache port
	CACHE_PORT_VALUE int = 6379
	// SCHEMA_NAME_VALUE defines schema name
	SCHEMA_NAME_VALUE string = "public"
	// ADMIN_USER_VALUE defines admin user
	ADMIN_USER_VALUE string = "root"
	// MISSING_COMMENT_VALUE defines missing comment value
	MISSING_COMMENT_VALUE int = 999
)

var (
	badSmtpPort int = SMTP_PORT_VALUE

	badSshPort int = SSH_PORT_VALUE

	badTelnetPort int = TELNET_PORT_VALUE

	badPoolMaxSize int = POOL_MAX_SIZE_VALUE

	badPoolMinSize int = POOL_MIN_SIZE_VALUE

	badConnectionTimeout int = CONNECTION_TIMEOUT_VALUE

	badServerVersion string = SERVER_VERSION_VALUE

	badBasePath string = BASE_PATH_VALUE

	badAuthToken string = AUTH_TOKEN_VALUE

	badAutoReload bool = AUTO_RELOAD_VALUE

	badStrictMode bool = STRICT_MODE_VALUE

	badLogEnabled bool = LOG_ENABLED_VALUE

	badAttemptsLimit int = ATTEMPTS_LIMIT_VALUE

	badWaitTimeMs int = WAIT_TIME_MS_VALUE

	badScaleFactor float64 = SCALE_FACTOR_VALUE

	badCacheHost string = CACHE_HOST_VALUE

	badCachePort int = CACHE_PORT_VALUE

	badSchemaName string = SCHEMA_NAME_VALUE

	badAdminUser string = ADMIN_USER_VALUE

	badMissingComment int = MISSING_COMMENT_VALUE
)
