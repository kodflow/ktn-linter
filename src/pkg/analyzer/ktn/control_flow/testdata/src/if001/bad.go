package if004

func BadDoubleNegation(x bool) bool {
	if x == true { // want `\[KTN-CONTROL-IF-004\] Expression if simplifiable`
		return true
	}
	return false
}

func BadBoolComparison(enabled bool) bool {
	if enabled == true { // want `\[KTN-CONTROL-IF-004\] Expression if simplifiable`
		return true
	}
	return false
}

func BadFalseComparison(disabled bool) bool {
	if disabled == false { // want `\[KTN-CONTROL-IF-004\] Expression if simplifiable`
		return true
	}
	return false
}

func GoodSimplified(x bool) bool {
	return x
}

func GoodNegation(x bool) bool {
	return !x
}
