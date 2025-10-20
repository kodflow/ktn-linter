package func010

// Good: 1 return value (unnamed OK)
func OneReturn() int {
	return 1
}

// Good: 2 return values (unnamed OK)
func TwoReturns() (int, error) {
	return 1, nil
}

// Good: 3 return values (unnamed OK - at limit)
func ThreeReturns() (int, string, error) {
	return 1, "test", nil
}

// Good: 4 return values with names
func FourNamedReturns() (count int, name string, valid bool, err error) {
	return 1, "test", true, nil
}

// Good: 5 return values with names
func FiveNamedReturns() (a int, b int, c string, d bool, e error) {
	return 1, 2, "test", true, nil
}

// Good: More than 3 returns but all named
func ManyNamedReturns() (id int, name string, age int, active bool, score float64) {
	return 1, "test", 30, true, 95.5
}

// Good: Function with no return value
func NoReturn() {
	x := 1
	_ = x
}

// Good: Test function with many unnamed returns (test functions are ignored)
func TestManyUnnamedReturns() (int, string, bool, error) {
	return 1, "test", true, nil
}

// Good: Benchmark function with many unnamed returns (benchmark functions are ignored)
func BenchmarkManyUnnamedReturns() (int, string, bool, error) {
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
	// Retour de la fonction
	return 1, 2, "point", nil
}

// TestFunction est une fonction de test (ignorée).
func TestSomething() {
	// Les fonctions de test ne sont pas vérifiées
	x := 1
	_ = x
}
