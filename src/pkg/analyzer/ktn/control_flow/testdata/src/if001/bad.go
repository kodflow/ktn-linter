package if004

func BadDoubleNegation(x bool) bool {
	// want `\[KTN-CONTROL-IF-004\] Expression if simplifiable`
	if x == true {
		return true
	}
	return false
}

func BadBoolComparison(enabled bool) bool {
	// want `\[KTN-CONTROL-IF-004\] Expression if simplifiable`
	if enabled == true {
		return true
	}
	return false
}

func BadFalseComparison(disabled bool) bool {
	// want `\[KTN-CONTROL-IF-004\] Expression if simplifiable`
	if disabled == false {
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
