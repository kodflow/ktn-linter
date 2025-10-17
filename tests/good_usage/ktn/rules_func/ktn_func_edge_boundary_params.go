package rules_func

// FiveParamsIsOk a exactement 5 paramètres (limite acceptable).
//
// Params:
//   - a: premier paramètre
//   - b: deuxième paramètre
//   - c: troisième paramètre
//   - d: quatrième paramètre
//   - e: cinquième paramètre
//
// Returns:
//   - int: résultat
func FiveParamsIsOk(a, b, c, d, e int) int {
	// Retourne la somme des cinq paramètres
	return a + b + c + d + e
}

// SixParamsWithConfig utilise une struct pour éviter trop de paramètres.
//
// Params:
//   - cfg: configuration avec tous les paramètres
//
// Returns:
//   - int: résultat
func SixParamsWithConfig(cfg SixParamsConfig) int {
	// Retourne la somme de tous les champs de la configuration
	return cfg.A + cfg.B + cfg.C + cfg.D + cfg.E + cfg.F
}

// SixParamsConfig contient les paramètres de configuration.
type SixParamsConfig struct {
	// A premier paramètre
	A int
	// B deuxième paramètre
	B int
	// C troisième paramètre
	C int
	// D quatrième paramètre
	D int
	// E cinquième paramètre
	E int
	// F sixième paramètre
	F int
}
