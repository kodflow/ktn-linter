package op001

// NOTE: L'analyzer est mal conçu - il détecte uniquement la division par le littéral 0,
// ce qui est déjà une erreur de compilation en Go.
// Il devrait plutôt détecter les divisions par des variables non vérifiées.
// TODO: Réimplémenter cet analyzer correctement.

func GoodDivisionWithCheck(x, y int) int {
	if y != 0 {
		return x / y
	}
	return 0
}
