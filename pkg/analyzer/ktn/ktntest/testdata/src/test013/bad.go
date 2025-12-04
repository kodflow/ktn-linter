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
		return "", errors.New("empty path")
	}
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
		return errors.New("input must be positive")
	}
	return nil
}

// GetVersion retourne la version.
// NE retourne PAS error → test table-driven avec erreur = over-engineering.
//
// Returns:
//   - string: version
func GetVersion() string {
	return "1.0.0"
}
