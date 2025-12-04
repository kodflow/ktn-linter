// Package test013 provides test cases for error cases.
package test013

import "errors"

// ParseConfig parse une configuration.
// Retourne error → le test DOIT avoir des cas d'erreur.
//
// Params:
//   - path: chemin du fichier
//
// Returns:
//   - string: contenu
//   - error: erreur si échec
func ParseConfig(path string) (string, error) {
	// Vérification chemin vide
	if path == "" {
		// Retour erreur
		return "", errors.New("empty path")
	}
	// Retour succès
	return "config:" + path, nil
}

// ValidateInput valide une entrée.
// Retourne error → le test DOIT avoir des cas d'erreur.
//
// Params:
//   - input: valeur à valider
//
// Returns:
//   - error: erreur si invalide
func ValidateInput(input int) error {
	// Vérification valeur positive
	if input <= 0 {
		// Retour erreur
		return errors.New("input must be positive")
	}
	// Retour succès
	return nil
}

// GetVersion retourne la version.
// NE retourne PAS error → test table-driven avec erreur = over-engineering.
//
// Returns:
//   - string: version
func GetVersion() string {
	// Retour version
	return "1.0.0"
}
