// Package func004 contient des exemples de fonctions privées utilisées correctement.
package func004

const (
	// MULTIPLIER est le facteur de multiplication utilisé dans compute.
	MULTIPLIER int = 2
)

// PublicFunction est une fonction publique.
//
// Returns:
//   - string: message
func PublicFunction() string {
	// Appel de la fonction privée
	if processData("test") {
		// Retour du helper
		return privateHelper()
	}
	// Retour vide
	return ""
}

// privateHelper est utilisée par PublicFunction.
//
// Returns:
//   - string: message
func privateHelper() string {
	// Retour du message
	return "helper"
}

// Calculator est une struct pour effectuer des calculs.
// Elle stocke une valeur interne et permet de la multiplier.
type Calculator struct {
	value int
}

// CalculatorInterface définit les méthodes publiques de Calculator.
type CalculatorInterface interface {
	Calculate() int
	Value() int
}

// NewCalculator crée une nouvelle instance de Calculator.
//
// Params:
//   - value: la valeur initiale
//
// Returns:
//   - *Calculator: nouvelle instance
func NewCalculator(value int) *Calculator {
	// Retour de la nouvelle instance
	return &Calculator{value: value}
}

// Calculate appelle la méthode privée.
//
// Returns:
//   - int: résultat
func (c *Calculator) Calculate() int {
	// Appel de la méthode privée
	return c.compute()
}

// compute est une méthode privée utilisée.
//
// Returns:
//   - int: résultat
func (c *Calculator) compute() int {
	// Retour de la valeur multipliée par la constante
	return c.value * MULTIPLIER
}

// Value retourne la valeur du calculateur.
//
// Returns:
//   - int: valeur actuelle
func (c *Calculator) Value() int {
	// Retour du champ value
	return c.value
}

// processData utilise validate en interne.
//
// Params:
//   - data: données
//
// Returns:
//   - bool: succès
func processData(data string) bool {
	// Appel de validate
	return validate(data)
}

// validate est utilisée par processData.
//
// Params:
//   - s: chaîne
//
// Returns:
//   - bool: valide
func validate(s string) bool {
	// Retour de la validation
	return len(s) > 0
}
