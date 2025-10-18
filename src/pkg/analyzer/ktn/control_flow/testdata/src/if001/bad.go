package if001

// TODO: L'analyseur IF-004 détecte les if/else qui retournent des booléens
// Il NE détecte PAS les comparaisons x == true (limitation actuelle)

// Ces exemples correspondent à ce que l'analyseur détecte
func BadIfElseTrue(x bool) bool {
	if x { // want `\[KTN-IF-004\].*`
		return true
	} else {
		return false
	}
}

func BadIfElseFalse(x bool) bool {
	if x { // want `\[KTN-IF-004\].*`
		return false
	} else {
		return true
	}
}

func GoodSimplified(x bool) bool {
	return x
}

func GoodNegation(x bool) bool {
	return !x
}
