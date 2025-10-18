// ════════════════════════════════════════════════════════════════════════════
// KTN-TEST-002: Fichier .go sans fichier _test.go correspondant
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//
//	Chaque fichier .go doit avoir son fichier _test.go correspondant.
//	Relation 1:1 obligatoire entre code et tests.
//
//	POURQUOI :
//	- Garantit que chaque fichier a des tests
//	- Facilite la navigation (même nom + _test)
//	- Maintient une couverture de test cohérente
//	- Évite les fichiers oubliés sans tests
//
// ❌ CAS INCORRECT : Ce fichier n'a pas de calculator_test.go
// ERREUR ATTENDUE: KTN-TEST-002 sur calculator.go
//
// ✅ CAS PARFAIT (voir target/) :
//
//	calculator.go + calculator_test.go (les deux existent)
//
// ════════════════════════════════════════════════════════════════════════════
package KTN_TEST_002

// Calculator effectue des calculs.
type Calculator struct{}

// Add additionne deux nombres.
//
// Params:
//   - a: premier nombre
//   - b: deuxième nombre
//
// Returns:
//   - int: résultat de l'addition
func (c *Calculator) Add(a int, b int) int {
	// Early return from function.
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
func (c *Calculator) Multiply(a int, b int) int {
	// Early return from function.
	return a * b
}

// NOTE: calculator_test.go n'existe PAS (violation KTN-TEST-002)
