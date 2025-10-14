package rules_var

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-001 : Variables groupées dans var ()
// ════════════════════════════════════════════════════════════════════════════

// Boolean configuration
// Ces variables configurent les fonctionnalités de l'application (mutables)
var (
	// EnableFeatureXV001Good active la fonctionnalité X
	EnableFeatureXV001Good bool = true
	// EnableDebugV001Good active le mode debug
	EnableDebugV001Good bool = false
	// isProductionV001Good indique si l'environnement est en production
	isProductionV001Good bool = true
)

// String configuration
// Ces variables configurent les thèmes de l'application
var (
	// ThemeAutoV001Good est l'identifiant du thème automatique
	ThemeAutoV001Good string = "auto"
	// ThemeCustomV001Good est l'identifiant du thème personnalisé
	ThemeCustomV001Good string = "custom"
)

// Integer configuration
// Ces variables configurent les limites entières (ajustables à runtime)
var (
	// MaxQueueSizeV001Good est la taille maximale de la queue
	MaxQueueSizeV001Good int16 = 10000
	// DefaultBufferSizeV001Good est la taille par défaut du buffer
	DefaultBufferSizeV001Good int16 = 4096
	// minCacheSizeV001Good est la taille minimale du cache
	minCacheSizeV001Good int16 = 512
)

// updateConfigV001Good modifie les configurations à runtime
func updateConfigV001Good() {
	EnableFeatureXV001Good = false
	EnableDebugV001Good = true
	isProductionV001Good = false
	ThemeAutoV001Good = "dark"
	ThemeCustomV001Good = "light"
	MaxQueueSizeV001Good = 20000
	DefaultBufferSizeV001Good = 8192
	minCacheSizeV001Good = 1024
}
