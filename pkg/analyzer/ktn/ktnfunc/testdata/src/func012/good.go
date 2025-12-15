// Package func012 provides good test cases.
package func012

const (
	// TWO_INT représente la valeur 2
	TWO_INT int = 2
	// THIRTY_INT représente la valeur 30
	THIRTY_INT int = 30
	// NINETY_FIVE_DOT représente la valeur 95.5
	NINETY_FIVE_DOT float64 = 95.5
)

// OneReturn retourne une seule valeur entière.
//
// Returns:
//   - int: valeur 1
func OneReturn() int {
	// Retourne la valeur 1
	return 1
}

// TwoReturns retourne deux valeurs.
//
// Returns:
//   - int: valeur 1
//   - error: nil
func TwoReturns() (int, error) {
	// Retourne un entier et nil pour l'erreur
	return 1, nil
}

// ThreeReturns retourne trois valeurs.
//
// Returns:
//   - int: valeur 1
//   - string: chaîne "test"
//   - error: nil
func ThreeReturns() (int, string, error) {
	// Retourne un entier, une chaîne et nil pour l'erreur
	return 1, "test", nil
}

// FourNamedReturns retourne quatre valeurs nommées.
//
// Returns:
//   - count: valeur 1
//   - name: chaîne "test"
//   - valid: true
//   - err: nil
func FourNamedReturns() (count int, name string, valid bool, err error) {
	// Retourne les quatre valeurs nommées
	return 1, "test", true, nil
}

// FiveNamedReturns retourne cinq valeurs nommées.
//
// Returns:
//   - a: valeur 1
//   - b: valeur 2
//   - c: chaîne "test"
//   - d: true
//   - e: nil
func FiveNamedReturns() (a int, b int, c string, d bool, e error) {
	// Retourne les cinq valeurs nommées
	return 1, TWO_INT, "test", true, nil
}

// ManyNamedReturns retourne plusieurs valeurs nommées.
//
// Returns:
//   - id: valeur 1
//   - name: chaîne "test"
//   - age: valeur 30
//   - active: true
//   - score: valeur 95.5
func ManyNamedReturns() (id int, name string, age int, active bool, score float64) {
	// Retourne les valeurs pour id, name, age, active et score
	return 1, "test", THIRTY_INT, true, NINETY_FIVE_DOT
}

// NoReturn ne retourne aucune valeur.
func NoReturn() {
	x := 1
	_ = x
}

// TestManyUnnamedReturns est une fonction de test avec plusieurs retours non nommés.
//
// Returns:
//   - int: valeur 1
//   - string: chaîne "test"
//   - bool: true
//   - error: nil
func TestManyUnnamedReturns() (int, string, bool, error) {
	// Retourne les valeurs de test
	return 1, "test", true, nil
}

// BenchmarkManyUnnamedReturns est une fonction de benchmark avec plusieurs retours non nommés.
//
// Returns:
//   - int: valeur 1
//   - string: chaîne "test"
//   - bool: true
//   - error: nil
func BenchmarkManyUnnamedReturns() (int, string, bool, error) {
	// Retourne les valeurs de benchmark
	return 1, "test", true, nil
}

// NoReturnValue ne retourne aucune valeur.
func NoReturnValue() {
	x := 1
	_ = x
}

// GetFourValuesCompact retourne quatre valeurs nommées (format compact).
//
// Returns:
//   - x, y: coordonnées
//   - name: nom
//   - err: erreur éventuelle
func GetFourValuesCompact() (x, y int, name string, err error) {
	// Retour des coordonnées et du nom
	return 1, TWO_INT, "point", nil
}

// testSomething est une fonction de test.
func testSomething() {
	// Les fonctions de test ne sont pas vérifiées
	x := 1
	_ = x
}

// init utilise les fonctions privées
func init() {
	// Appel de testSomething
	testSomething()
}
