package test002

import "errors"

const (
	// MIN_LENGTH longueur minimale pour les strings
	MIN_LENGTH int = 1
	// MAX_VALUE valeur maximale autorisée
	MAX_VALUE int = 100
)

// Validate valide une valeur numérique.
//
// Params:
//   - value: valeur à valider
//
// Returns:
//   - error: erreur si invalide, nil sinon
func Validate(value int) error {
	// Vérification valeur négative
	if value < 0 {
		// Retour erreur
		return errors.New("value must be non-negative")
	}
	// Vérification valeur trop grande
	if value > MAX_VALUE {
		// Retour erreur
		return errors.New("value exceeds maximum")
	}
	// Retour succès
	return nil
}

// Transform transforme une string.
//
// Params:
//   - input: string à transformer
//
// Returns:
//   - string: string transformée
//   - error: erreur si vide
func Transform(input string) (string, error) {
	// Vérification longueur
	if len(input) < MIN_LENGTH {
		// Retour erreur
		return "", errors.New("input too short")
	}
	// Retour résultat
	return "[" + input + "]", nil
}
