// Package test002 provides test cases for error scenarios.
package test002

import "errors"

const (
	// EMPTY_LENGTH longueur vide
	EMPTY_LENGTH int = 0
)

// CheckValue vérifie une valeur.
//
// Params:
//   - val: valeur à vérifier
//
// Returns:
//   - error: erreur si invalide
func CheckValue(val int) error {
	// Vérification valeur nulle
	if val == EMPTY_LENGTH {
		// Retour erreur
		return errors.New("value cannot be zero")
	}
	// Retour succès
	return nil
}

// Format formatte une string.
//
// Params:
//   - s: string à formater
//
// Returns:
//   - string: résultat
//   - error: erreur si problème
func Format(s string) (string, error) {
	// Vérification longueur
	if len(s) == EMPTY_LENGTH {
		// Retour erreur
		return "", errors.New("cannot format empty string")
	}
	// Retour résultat
	return ">" + s + "<", nil
}
