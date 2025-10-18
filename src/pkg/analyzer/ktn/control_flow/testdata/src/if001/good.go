package if001

// Cas corrects - l'analyseur ne devrait pas rapporter d'erreur

// GoodExample - fonction simple
func GoodExample() {
	// Code conforme aux règles
}

// GoodFunction - calcul simple
func GoodFunction(x int) int {
	return x * 2
}

// GoodWithInit - if avec initialisation (ne devrait PAS être détecté)
func GoodWithInit(x int) bool {
	if y := x > 0; y {
		return true
	} else {
		return false
	}
}

// GoodMultipleStatementsBody - if avec plusieurs statements dans body
func GoodMultipleStatementsBody(x bool) bool {
	if x {
		println("debug")
		return true
	} else {
		return false
	}
}

// GoodMultipleStatementsElse - if avec plusieurs statements dans else
func GoodMultipleStatementsElse(x bool) bool {
	if x {
		return true
	} else {
		println("debug")
		return false
	}
}

// GoodNonBooleanReturn - if qui retourne des non-booléens
func GoodNonBooleanReturn(x bool) int {
	if x {
		return 1
	} else {
		return 0
	}
}

// GoodMultipleReturnValues - if qui retourne plusieurs valeurs
func GoodMultipleReturnValues(x bool) (bool, error) {
	if x {
		return true, nil
	} else {
		return false, nil
	}
}

// GoodSameValue - if/else qui retournent la même valeur
func GoodSameValue(x bool) bool {
	if x {
		return true
	} else {
		return true
	}
}

// GoodNoElse - if sans else
func GoodNoElse(x bool) bool {
	if x {
		return true
	}
	return false
}

// GoodNonLiteralReturn - if qui retourne des variables
func GoodNonLiteralReturn(x bool, a bool, b bool) bool {
	if x {
		return a
	} else {
		return b
	}
}
