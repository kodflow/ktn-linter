// Package test012 is used for testing passthrough test detection.
package test012

// ProcessData traite des données.
//
// Params:
//   - data: données à traiter
//
// Returns:
//   - string: résultat
func ProcessData(data string) string {
	// Traitement des données
	return "processed:" + data
}

// GetCount retourne un compteur.
//
// Returns:
//   - int: valeur du compteur
func GetCount() int {
	// Retour de la valeur
	return 42
}

// IsValid vérifie si une valeur est valide.
//
// Params:
//   - value: valeur à vérifier
//
// Returns:
//   - bool: true si valide
func IsValid(value int) bool {
	// Vérification de la valeur
	return value > 0
}

// GetError retourne une erreur pour les tests.
//
// Params:
//   - fail: si true, retourne une erreur
//
// Returns:
//   - error: erreur si fail est true
func GetError(fail bool) error {
	// Pas d'erreur pour les tests
	return nil
}
