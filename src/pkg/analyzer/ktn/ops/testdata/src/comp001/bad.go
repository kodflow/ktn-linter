package comp001

func BadRedundantTrue(x bool) bool {
	// want `\[KTN-OPS-COMP-001\] Comparaison redondante`
	if x == true {
		return true
	}
	return false
}

func BadRedundantFalse(x bool) bool {
	// want `\[KTN-OPS-COMP-001\] Comparaison redondante`
	if x == false {
		return false
	}
	return true
}

func GoodDirectBool(x bool) bool {
	return x
}

func GoodNegation(x bool) bool {
	return !x
}
