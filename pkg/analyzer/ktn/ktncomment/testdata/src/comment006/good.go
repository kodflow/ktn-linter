// Package comment006 provides good test cases.
package comment006

// GoodNoParams effectue une action simple sans paramètres.
//
// Returns:
//   - string: le résultat de l'opération
func GoodNoParams() string {
	// Retourne le résultat de l'opération
	return "result"
}

// GoodWithParams traite des données avec des paramètres.
//
// Params:
//   - _data: les données à traiter
//   - _count: le nombre d'itérations
//
// Returns:
//   - bool: true si le traitement a réussi
//   - error: une erreur si le traitement a échoué
func GoodWithParams(_data string, _count int) (bool, error) {
	// Retourne succès du traitement sans erreur
	return true, nil
}

// GoodNoReturn effectue une action sans retourner de valeur.
//
// Params:
//   - _msg: le message à afficher (non utilisé)
func GoodNoReturn(_msg string) {
	// Fonction sans retour
}

// GoodNoParamsNoReturn effectue une opération sans paramètres ni retour.
func GoodNoParamsNoReturn() {
	// Do something
}

// GoodComplete démontre une documentation complète avec exemple.
//
// Params:
//   - input: les données d'entrée
//
// Returns:
//   - string: les données formatées
//
// Example:
//
//	result := GoodComplete("test")
func GoodComplete(input string) string {
	// Retourne les données d'entrée formatées
	return input
}

// GoodVariadic calcule la somme de valeurs variadiques.
//
// Params:
//   - values: liste de valeurs à sommer
//
// Returns:
//   - int: la somme totale
func GoodVariadic(values ...int) int {
	// Calcul de la somme
	sum := 0
	// Itération sur les valeurs
	for _, v := range values {
		// Addition de la valeur
		sum += v
	}
	// Retourne la somme
	return sum
}

// GoodMultipleReturns retourne plusieurs valeurs.
//
// Params:
//   - x: première valeur
//   - y: deuxième valeur
//
// Returns:
//   - int: la somme
//   - int: la différence
//   - error: erreur éventuelle
func GoodMultipleReturns(x, y int) (int, int, error) {
	// Retourne somme, différence et nil
	return x + y, x - y, nil
}

// GoodNamedReturns utilise des retours nommés.
//
// Params:
//   - input: valeur d'entrée
//
// Returns:
//   - result: le résultat calculé
//   - err: erreur éventuelle
func GoodNamedReturns(input int) (result int, err error) {
	// Calcul du résultat
	result = input + input
	// Retourne les valeurs nommées
	return result, nil
}

// GoodUnderscoreParams utilise des paramètres underscore.
//
// Params:
//   - _unused: paramètre non utilisé
//
// Returns:
//   - int: valeur de retour
func GoodUnderscoreParams(_unused string) int {
	// Retourne zéro
	return 0
}

// init initialise le package.
func init() {
	// Appel des fonctions pour éviter dead code
	_ = GoodNoParams()
	_, _ = GoodWithParams("data", 1)
	GoodNoReturn("msg")
	GoodNoParamsNoReturn()
	_ = GoodComplete("input")
	_ = GoodVariadic(1, 1, 1)
	_, _, _ = GoodMultipleReturns(1, 1)
	_, _ = GoodNamedReturns(1)
	_ = GoodUnderscoreParams("unused")
}
