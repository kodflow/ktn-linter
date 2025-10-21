package var004

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
	badSmtpPort int = SMTP_PORT_VALUE // want "KTN-VAR-004"

	badSshPort int = SSH_PORT_VALUE // want "KTN-VAR-004"

	badTelnetPort int = TELNET_PORT_VALUE // want "KTN-VAR-004"

	badPoolMaxSize int = POOL_MAX_SIZE_VALUE // want "KTN-VAR-004"

	badPoolMinSize int = POOL_MIN_SIZE_VALUE // want "KTN-VAR-004"

	badConnectionTimeout int = CONNECTION_TIMEOUT_VALUE // want "KTN-VAR-004"

	badServerVersion string = SERVER_VERSION_VALUE // want "KTN-VAR-004"

	badBasePath string = BASE_PATH_VALUE // want "KTN-VAR-004"

	badAuthToken string = AUTH_TOKEN_VALUE // want "KTN-VAR-004"

	badAutoReload bool = AUTO_RELOAD_VALUE // want "KTN-VAR-004"

	badStrictMode bool = STRICT_MODE_VALUE // want "KTN-VAR-004"

	badLogEnabled bool = LOG_ENABLED_VALUE // want "KTN-VAR-004"

	badAttemptsLimit int = ATTEMPTS_LIMIT_VALUE // want "KTN-VAR-004"

	badWaitTimeMs int = WAIT_TIME_MS_VALUE // want "KTN-VAR-004"

	badScaleFactor float64 = SCALE_FACTOR_VALUE // want "KTN-VAR-004"

	badCacheHost string = CACHE_HOST_VALUE // want "KTN-VAR-004"

	badCachePort int = CACHE_PORT_VALUE // want "KTN-VAR-004"

	badSchemaName string = SCHEMA_NAME_VALUE // want "KTN-VAR-004"

	badAdminUser string = ADMIN_USER_VALUE // want "KTN-VAR-004"

	badMissingComment int = MISSING_COMMENT_VALUE // want "KTN-VAR-004"
)
