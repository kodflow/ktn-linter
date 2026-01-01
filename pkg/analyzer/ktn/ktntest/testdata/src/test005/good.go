// Package test005 provides simple test utilities.
package test005

const (
	// answerValue valeur de réponse
	answerValue int = 42
)

// GoodFunction est une fonction qui sera testée.
//
// Returns:
//   - int: valeur de retour
func GoodFunction() int {
	// Retour de la valeur
	return answerValue
}
