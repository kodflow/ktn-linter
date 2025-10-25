package test006

import "errors"

const (
	// MIN_VALUE valeur minimale
	MIN_VALUE int = 0
	// MAX_LENGTH longueur maximale
	MAX_LENGTH int = 100
)

// processData traite des données (fonction privée).
//
// Params:
//   - data: données à traiter
//
// Returns:
//   - string: résultat
//   - error: erreur si invalide
func processData(data string) (string, error) {
	// Vérification longueur
	if len(data) > MAX_LENGTH {
		// Retour erreur
		return "", errors.New("data too long")
	}
	// Retour résultat
	return "processed:" + data, nil
}

// validateNumber valide un nombre (fonction privée).
//
// Params:
//   - num: nombre à valider
//
// Returns:
//   - error: erreur si invalide
func validateNumber(num int) error {
	// Vérification valeur négative
	if num < MIN_VALUE {
		// Retour erreur
		return errors.New("number must be non-negative")
	}
	// Retour succès
	return nil
}
