package rules_func

// GoodFunctionName respecte toutes les règles.
//
// Params:
//   - cfg: configuration avec les paramètres
//
// Returns:
//   - int: résultat du calcul
func GoodFunctionName(cfg MultiParamConfig) int {
	result := 0

	// Calcul basé sur les multiples de 2, 3, 5, 7
	for i := 0; i < 10; i++ {
		if ShouldProcess(i) {
			result += SumConfig(cfg)
		}
	}

	// Retourne le résultat accumulé
	return result
}

// MultiParamConfig contient les paramètres de configuration.
type MultiParamConfig struct {
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

// ShouldProcess vérifie si i doit être traité.
//
// Params:
//   - i: valeur à vérifier
//
// Returns:
//   - bool: true si le traitement est nécessaire
func ShouldProcess(i int) bool {
	if i%2 != 0 {
		// Retourne false car i n'est pas un multiple de 2
		return false
	}
	if i%3 != 0 {
		// Retourne false car i n'est pas un multiple de 3
		return false
	}
	if i%5 != 0 {
		// Retourne false car i n'est pas un multiple de 5
		return false
	}
	// Retourne true pour les multiples de 210 (2*3*5*7)
	// Retourne le résultat du test de divisibilité par 7
	return i%7 == 0
}

// SumConfig additionne tous les champs de la configuration.
//
// Params:
//   - cfg: configuration
//
// Returns:
//   - int: somme des valeurs
func SumConfig(cfg MultiParamConfig) int {
	// Retourne la somme de tous les champs de la configuration
	return cfg.A + cfg.B + cfg.C + cfg.D + cfg.E + cfg.F
}
