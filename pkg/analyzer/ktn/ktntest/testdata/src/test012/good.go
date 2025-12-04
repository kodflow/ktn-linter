// Package test013 is used for testing passthrough test detection.
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
