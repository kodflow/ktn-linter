// ════════════════════════════════════════════════════════════════════════════
// KTN-INTERFACE-006: Violations avec suffixes -er
// ════════════════════════════════════════════════════════════════════════════

package KTN_INTERFACE_006

// CacheManager interface sans suffixe -er (violation)
type CacheManager interface {
	Get(key string) (string, bool)
	Set(key string, value string)
}

// Logger interface sans suffixe -er (violation)
type Logger interface {
	Info(msg string)
	Error(msg string)
}

// marker interface vide sans documentation (violation)
type marker interface{}
