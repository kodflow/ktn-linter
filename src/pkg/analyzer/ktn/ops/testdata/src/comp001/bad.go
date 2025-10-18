package comp001

func BadRedundantTrue(x bool) bool {
	if x == true { // want `\[KTN-COMP-001\].*`
		return true
	}
	return false
}

func BadRedundantFalse(x bool) bool {
	if x == false { // want `\[KTN-COMP-001\].*`
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
