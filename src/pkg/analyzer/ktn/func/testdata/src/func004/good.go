package func004

// GoodWithReturnsSection calcule une valeur.
//
// Returns:
//   - int: le résultat du calcul
//   - error: erreur éventuelle
func GoodWithReturnsSection() (int, error) {
	return 42, nil
}

// GoodSingleReturn retourne un résultat.
//
// Returns:
//   - string: le message de succès
func GoodSingleReturn() string {
	return "success"
}

// GoodNoReturn ne retourne rien - pas besoin de section Returns.
func GoodNoReturn() {
	// fonction sans retour
}

// GoodMultipleReturns effectue un calcul complexe.
//
// Returns:
//   - float64: le total calculé
//   - int: le nombre d'éléments
//   - error: erreur si calcul impossible
func GoodMultipleReturns() (float64, int, error) {
	return 3.14, 10, nil
}
