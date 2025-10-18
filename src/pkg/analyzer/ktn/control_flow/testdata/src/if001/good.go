package if001

// Cas corrects - early returns sans else (pattern idiomatique Go)

// GoodExample - fonction simple
func GoodExample() {
	// Code conforme aux règles
}

// GoodFunction - calcul simple
func GoodFunction(x int) int {
	return x * 2
}

// GoodWithInit - if avec initialisation
func GoodWithInit(x int) bool {
	if y := x > 0; y {
		return true
	}
	return false
}

// GoodMultipleStatementsBody - if avec plusieurs statements
func GoodMultipleStatementsBody(x bool) bool {
	if x {
		println("debug")
		return true
	}
	return false
}

// GoodMultipleStatementsAfter - plusieurs statements après if
func GoodMultipleStatementsAfter(x bool) bool {
	if x {
		return true
	}
	println("debug")
	return false
}

// GoodNonBooleanReturn - if qui retourne des non-booléens
func GoodNonBooleanReturn(x bool) int {
	if x {
		return 1
	}
	return 0
}

// GoodMultipleReturnValues - if qui retourne plusieurs valeurs
func GoodMultipleReturnValues(x bool) (bool, error) {
	if x {
		return true, nil
	}
	return false, nil
}

// GoodSameValue - if qui retourne la même valeur
func GoodSameValue(x bool) bool {
	if x {
		return true
	}
	return true
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
	}
	return b
}

// GoodNestedChecks - checks imbriqués sans else
func GoodNestedChecks(x, y int) int {
	if x > 0 {
		if y > 0 {
			return x + y
		}
		return x
	}
	return 0
}

// GoodComplexCondition - condition complexe sans else
func GoodComplexCondition(a, b int) bool {
	if a > 0 && b > 0 {
		return true
	}
	return false
}
