// Package test011 provides data processing utilities.
package test011

import "errors"

const (
	// countValue valeur du compteur
	countValue int = 42
)

// ProcessData traite des données.
// Retourne error → le test a des cas d'erreur = OK.
//
// Params:
//   - data: données à traiter
//
// Returns:
//   - string: résultat
//   - error: erreur si échec
func ProcessData(data string) (string, error) {
	// Vérification données vides
	if data == "" {
		// Retour erreur
		return "", errors.New("empty data")
	}
	// Retour du résultat
	return "processed:" + data, nil
}

// GetName retourne un nom.
// NE retourne PAS error → test simple sans cas d'erreur = OK.
//
// Returns:
//   - string: nom
func GetName() string {
	// Retour du nom
	return "test"
}

// GetCount retourne un compteur.
// NE retourne PAS error → test table-driven SANS cas d'erreur = OK.
//
// Returns:
//   - int: compteur
func GetCount() int {
	// Retour du compteur
	return countValue
}
