// Good examples for the func007 test case.
package comment006

// GoodNoParams effectue une action simple sans paramètres
//
// Returns:
//   - string: le résultat de l'opération
func GoodNoParams() string {
	// Retourne le résultat de l'opération
	return "result"
}

// GoodWithParams traite des données avec des paramètres
//
// Params:
//   - _data: les données à traiter (non utilisées dans cet exemple)
//   - _count: le nombre d'itérations (non utilisé dans cet exemple)
//
// Returns:
//   - bool: true si le traitement a réussi
//   - error: une erreur si le traitement a échoué
func GoodWithParams(_data string, _count int) (bool, error) {
	// Retourne succès du traitement sans erreur
	return true, nil
}

// GoodNoReturn effectue une action sans retourner de valeur
//
// Params:
//   - _msg: le message à afficher (non utilisé dans cet exemple)
func GoodNoReturn(_msg string) {
	// Do something
}

// GoodComplete démontre une documentation complète avec exemple
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
