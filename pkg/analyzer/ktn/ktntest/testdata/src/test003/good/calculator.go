package test003

const (
	// ZERO_DIVISOR valeur du diviseur zéro
	ZERO_DIVISOR int = 0
)

// Add additionne deux nombres.
//
// Params:
//   - a: premier nombre
//   - b: deuxième nombre
//
// Returns:
//   - int: somme de a et b
func Add(a, b int) int {
	// Retour de la somme
	return a + b
}

// Subtract soustrait deux nombres.
//
// Params:
//   - a: premier nombre
//   - b: deuxième nombre
//
// Returns:
//   - int: différence de a et b
func Subtract(a, b int) int {
	// Retour de la différence
	return a - b
}

// Multiply multiplie deux nombres.
//
// Params:
//   - a: premier nombre
//   - b: deuxième nombre
//
// Returns:
//   - int: produit de a et b
func Multiply(a, b int) int {
	// Retour du produit
	return a * b
}

// Divide divise deux nombres.
//
// Params:
//   - a: dividende
//   - b: diviseur
//
// Returns:
//   - int: quotient de a divisé par b
//   - bool: false si division par zéro, true sinon
func Divide(a, b int) (int, bool) {
	// Vérification division par zéro
	if b == ZERO_DIVISOR {
		// Retour erreur
		return 0, false
	}
	// Retour du quotient
	return a / b, true
}
