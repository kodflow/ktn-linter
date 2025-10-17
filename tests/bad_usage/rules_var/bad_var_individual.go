package rules_var

// ANTI-PATTERN: Variables déclarées individuellement
// Viole KTN-VAR-001

// GlobalConfig variable individuelle - MAUVAIS !
var GlobalConfig string

// DatabaseURL variable individuelle
var DatabaseURL string

// MaxConnections variable individuelle
var MaxConnections int

// EnableDebug variable individuelle
var EnableDebug bool

// DefaultTimeout variable individuelle
var DefaultTimeout int

// CacheTTL variable individuelle
var CacheTTL int

// APIKey variable individuelle
var APIKey string

// SecretKey variable individuelle
var SecretKey string

// ServerPort variable individuelle
var ServerPort int

// HostName variable individuelle
var HostName string
