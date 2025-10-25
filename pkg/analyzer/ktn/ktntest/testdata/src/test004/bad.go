package test004

import "errors"

const (
	// ZERO_VALUE valeur zéro
	ZERO_VALUE int = 0
	// POSITIVE_THRESHOLD seuil positif
	POSITIVE_THRESHOLD int = 0
)

// CheckPositive vérifie si un nombre est positif.
//
// Params:
//   - n: nombre à vérifier
//
// Returns:
//   - error: erreur si négatif ou zéro, nil sinon
func CheckPositive(n int) error {
	// Vérification nombre positif
	if n <= POSITIVE_THRESHOLD {
		// Retour erreur
		return errors.New("number must be positive")
	}
	// Retour succès
	return nil
}

// FormatString formate une string.
//
// Params:
//   - s: string à formater
//
// Returns:
//   - string: string formatée
//   - error: erreur si vide
func FormatString(s string) (string, error) {
	// Vérification string vide
	if s == "" {
		// Retour erreur
		return "", errors.New("empty string")
	}
	// Retour résultat
	return "[" + s + "]", nil
}
