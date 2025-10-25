package test004

import "errors"

const (
	// ZERO_VALUE valeur zéro
	ZERO_VALUE int = 0
	// POSITIVE_THRESHOLD seuil positif
	POSITIVE_THRESHOLD int = 0
)

// ValidateInput valide une entrée utilisateur.
//
// Params:
//   - input: valeur à valider
//
// Returns:
//   - error: erreur si invalide, nil sinon
func ValidateInput(input int) error {
	// Vérification valeur positive
	if input <= POSITIVE_THRESHOLD {
		// Retour erreur
		return errors.New("input must be positive")
	}
	// Retour succès
	return nil
}

// ProcessData traite des données.
//
// Params:
//   - data: données à traiter
//
// Returns:
//   - string: résultat du traitement
//   - error: erreur si échec
func ProcessData(data string) (string, error) {
	// Vérification données vides
	if data == "" {
		// Retour erreur
		return "", errors.New("empty data")
	}
	// Retour résultat
	return "processed: " + data, nil
}
