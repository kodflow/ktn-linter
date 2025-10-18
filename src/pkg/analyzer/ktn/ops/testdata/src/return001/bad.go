package return004

func BadNakedReturn(x int) (result int) {
	result = x * 2
	// want `\[KTN-OPS-RETURN-004\] Naked return détecté`
	return
}

func BadNakedReturnMultiple(a, b int) (sum int, product int) {
	sum = a + b
	product = a * b
	// want `\[KTN-OPS-RETURN-004\] Naked return détecté`
	return
}

func GoodExplicitReturn(x int) (result int) {
	result = x * 2
	return result
}

func GoodSimpleReturn(x int) int {
	return x * 2
}
