package rules_var

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-008 : Nom avec MixedCaps (pas d'underscore)
// ════════════════════════════════════════════════════════════════════════════

// HTTP status codes
// Ces variables représentent les codes de statut HTTP standards
var (
	// HTTPOKV008Good représente le code HTTP 200
	HTTPOKV008Good int = 200
	// HTTPNotFoundV008Good représente le code HTTP 404
	HTTPNotFoundV008Good int = 404
)

// updateStatusCodesV008Good modifie les codes HTTP à runtime
func updateStatusCodesV008Good() {
	HTTPOKV008Good = 201
	HTTPNotFoundV008Good = 410
}
