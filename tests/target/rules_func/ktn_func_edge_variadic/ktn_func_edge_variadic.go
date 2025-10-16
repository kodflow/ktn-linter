package rules_func

import "fmt"

// Fonctions variadiques correctement documentées.

// Sum calcule la somme d'un nombre variable d'entiers.
//
// Params:
//   - nums: liste variable de nombres à additionner
//
// Returns:
//   - int: la somme de tous les nombres
func Sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	// Retourne la somme totale
	return total
}

// PrintfWrapper enveloppe fmt.Printf avec un formatage personnalisé.
//
// Params:
//   - format: chaîne de format Printf
//   - args: arguments variables pour le formatage
func PrintfWrapper(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// ProcessItems préfixe chaque élément d'une liste variable.
//
// Params:
//   - prefix: préfixe à ajouter à chaque élément
//   - items: liste variable d'éléments à préfixer
//
// Returns:
//   - []string: liste des éléments préfixés
func ProcessItems(prefix string, items ...string) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = prefix + item
	}
	// Retourne la liste des éléments préfixés
	return result
}

// MergeAndProcess fusionne et traite des valeurs (utilise un struct pour éviter trop de params).
//
// Params:
//   - config: configuration des valeurs de base
//   - extra: valeurs supplémentaires variables
//
// Returns:
//   - int: somme de toutes les valeurs
func MergeAndProcess(config MergeConfig, extra ...int) int {
	total := config.A + config.B + config.C + config.D
	for _, e := range extra {
		total += e
	}
	// Retourne la somme totale
	return total
}

// MergeConfig configuration pour MergeAndProcess.
type MergeConfig struct {
	// A est la première valeur de base
	A int
	// B est la deuxième valeur de base
	B int
	// C est la troisième valeur de base
	C int
	// D est la quatrième valeur de base
	D int
}

// ApplyMultiplier applique un multiplicateur à des valeurs variables.
//
// Params:
//   - multiplier: facteur de multiplication
//   - values: liste variable de valeurs à multiplier
//
// Returns:
//   - []float64: liste des valeurs multipliées
func ApplyMultiplier(multiplier int, values ...float64) []float64 {
	result := make([]float64, 0, len(values))

	// Filtrage et transformation des valeurs
	for _, v := range values {
		if processed, ok := processValueWithMultiplier(v, multiplier); ok {
			result = append(result, processed)
		}
	}

	// Retourne la liste des valeurs multipliées
	return result
}

// processValueWithMultiplier traite une valeur avec le multiplicateur.
//
// Params:
//   - v: la valeur à traiter
//   - multiplier: facteur de multiplication
//
// Returns:
//   - float64: la valeur traitée
//   - bool: true si la valeur doit être incluse
func processValueWithMultiplier(v float64, multiplier int) (float64, bool) {
	// Filtrer les valeurs négatives ou nulles
	if v <= 0 || v >= 100 {
		// Retourne 0 et false pour exclure la valeur
		return 0, false
	}

	// Application du multiplicateur si > 1
	if multiplier > 1 {
		// Retourne la valeur multipliée et true
		return v * float64(multiplier), true
	}
	// Retourne la valeur inchangée et true
	return v, true
}
