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

// BadComplexCondition - condition complexe mais pattern simplifiable
func BadComplexCondition(a, b bool) bool {
	if a && b { // want `\[KTN-IF-004\].*`
		return true
	} else {
		return false
	}
}

// Cas qui ne doivent PAS être détectés

func GoodWithInitBad(x int) bool {
	// if avec initialisation - ne devrait PAS être signalé
	if y := x > 0; y {
		return true
	}
	return false
}

func GoodMultipleStatementsBad(x bool) bool {
	// if avec plusieurs statements - ne devrait PAS être signalé
	if x {
		println("debug")
		return true
	}
	return false
}

func GoodNonBoolReturnBad(x bool) int {
	// if qui ne retourne pas un booléen littéral
	if x {
		return 1
	} else {
		return 0
	}
}

func GoodSameBoolValueBad(x bool) bool {
	// if/else qui retournent la même valeur (pas simplifiable)
	if x {
		return true
	} else {
		return true
	}
}
