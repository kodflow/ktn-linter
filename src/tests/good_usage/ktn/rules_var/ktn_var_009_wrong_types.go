package rules_var

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-009 : Nom en MixedCaps (pas ALL_CAPS)
// ════════════════════════════════════════════════════════════════════════════

// Buffer configuration
// Ces variables configurent les buffers et timeouts
var (
	// MaxBufferSizeV009Good est la taille maximale du buffer
	MaxBufferSizeV009Good int = 1024
	// DefaultTimeoutSecondsV009Good est le timeout par défaut en secondes
	DefaultTimeoutSecondsV009Good int = 30
)

// Theme configuration - Toutes les variables du même thème regroupées
// Ces variables définissent les thèmes disponibles (tous configurables ensemble)
var (
	// ThemeLightV009Good est l'identifiant du thème clair
	ThemeLightV009Good string = "light"
	// ThemeDarkV009Good est l'identifiant du thème sombre
	ThemeDarkV009Good string = "dark"
	// ThemeHighContrastV009Good est l'identifiant du thème à haut contraste
	ThemeHighContrastV009Good string = "high-contrast"
	// ThemeSepiaV009Good est l'identifiant du thème sépia
	ThemeSepiaV009Good string = "sepia"
)

// updateSettingsV009Good modifie les configurations à runtime
func updateSettingsV009Good() {
	MaxBufferSizeV009Good = 2048
	DefaultTimeoutSecondsV009Good = 60
	ThemeLightV009Good = "bright"
	ThemeDarkV009Good = "night"
	ThemeHighContrastV009Good = "accessibility"
	ThemeSepiaV009Good = "vintage"
}
