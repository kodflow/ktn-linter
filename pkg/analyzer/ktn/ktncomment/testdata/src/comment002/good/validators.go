// Package good provides validators for user input and system configuration.
package good

// ValidateInput valide les données d'entrée.
//
// Params:
//   - input: données à valider
//
// Returns:
//   - bool: true si valide
func ValidateInput(input string) bool {
	// Retour du résultat de validation
	return len(input) > 0
}
