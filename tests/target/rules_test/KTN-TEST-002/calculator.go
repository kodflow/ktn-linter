// ════════════════════════════════════════════════════════════════════════════
// KTN-TEST-002: Fichier .go avec son _test.go correspondant (✅ CORRIGÉ)
// ════════════════════════════════════════════════════════════════════════════
package KTN_TEST_002_GOOD

// CalculatorData effectue des calculs.
type CalculatorData struct{}

// Add additionne deux nombres.
//
// Params:
//   - a: premier nombre
//   - b: deuxième nombre
//
// Returns:
//   - int: résultat de l'addition
func (c *CalculatorData) Add(a int, b int) int {
	return a + b
}

// Multiply multiplie deux nombres.
//
// Params:
//   - a: premier nombre
//   - b: deuxième nombre
//
// Returns:
//   - int: résultat de la multiplication
func (c *CalculatorData) Multiply(a int, b int) int {
	return a * b
}

// ✅ CORRIGÉ: calculator_test.go existe maintenant
