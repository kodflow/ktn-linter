// Bad examples for the var011 test case.
package var011

const (
	// LOOP_MAX_ITERATIONS est le nombre maximum d'itérations
	LOOP_MAX_ITERATIONS int = 10
	// MULTIPLIER_VALUE est le multiplicateur utilisé
	MULTIPLIER_VALUE int = 2
)

// badShadowingCount démontre le shadowing d'une variable non exemptée.
//
// Note: Les variables err, ok, ctx sont exemptées car ce sont des patterns idiomatiques Go.
func badShadowingCount() {
	count := 0
	// Boucle sur les itérations
	for range LOOP_MAX_ITERATIONS {
		count := count * MULTIPLIER_VALUE // want "KTN-VAR-011: shadowing de la variable 'count'"
		_ = count
	}
	_ = count
}

// badShadowingValue démontre le shadowing d'une variable value.
func badShadowingValue() {
	value := "outer"
	// Bloc imbriqué
	{
		value := "inner" // want "KTN-VAR-011: shadowing de la variable 'value'"
		_ = value
	}
	_ = value
}

// badShadowingResult démontre le shadowing d'une variable result.
func badShadowingResult() {
	result, _ := doSomething()
	// Bloc if
	if result > 0 {
		result := result * MULTIPLIER_VALUE // want "KTN-VAR-011: shadowing de la variable 'result'"
		_ = result
	}
	_ = result
}

// badShadowingData démontre le shadowing d'une variable data.
func badShadowingData() {
	data := []int{1, 2, 3}
	// Boucle range
	for i := range data {
		data := append(data, i) // want "KTN-VAR-011: shadowing de la variable 'data'"
		_ = data
	}
	_ = data
}

// badShadowingName démontre le shadowing dans une fonction.
//
// Params:
//   - name: nom passé en paramètre
func badShadowingName(name string) {
	// Bloc imbriqué
	{
		name := "shadowed" // want "KTN-VAR-011: shadowing de la variable 'name'"
		_ = name
	}
	_ = name
}

// doSomething effectue une opération.
//
// Returns:
//   - int: résultat de l'opération
//   - error: erreur éventuelle
func doSomething() (int, error) {
	// Retour avec résultat
	return 0, nil
}

// init utilise les fonctions privées
func init() {
	// Appel de badShadowingCount
	badShadowingCount()
	// Appel de badShadowingValue
	badShadowingValue()
	// Appel de badShadowingResult
	badShadowingResult()
	// Appel de badShadowingData
	badShadowingData()
	// Appel de badShadowingName
	badShadowingName("test")
	// Appel de doSomething
	doSomething()
}
